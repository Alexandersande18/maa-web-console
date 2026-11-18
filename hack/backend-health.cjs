// Maintenance-only dependencies are supplied by `npm ci` in web/.
const fs = require('fs');
const path = require('path');
const http = require('http');
const https = require('https');
const yaml = require('../web/node_modules/js-yaml');
const WebSocket = require('../web/node_modules/ws');
const root = path.resolve(__dirname, '..');
const args = process.argv.slice(2);
const mode = args[0];
const option = (name, fallback) => args.includes(name) ? args[args.indexOf(name) + 1] : fallback;
const base = new URL(option('--base-url', 'http://127.0.0.1:5244'));
const snapshotPath = option('--snapshot');
let stage = 'read-configuration';
function request(method, route, body, token) {
  return new Promise((resolve, reject) => {
    const data = body === undefined ? undefined : JSON.stringify(body);
    const url = new URL(route, base);
    const transport = url.protocol === 'https:' ? https : http;
    const req = transport.request(url, { method, headers: {
      ...(data ? {'Content-Type':'application/json', 'Content-Length':Buffer.byteLength(data)} : {}),
      ...(token ? {Authorization:token} : {})
    }}, res => {
      let text = '';
      res.setEncoding('utf8');
      res.on('data', b => { text += b; });
      res.on('end', () => {
        try {
          const value = JSON.parse(text);
          if (res.statusCode !== 200 || value.code !== 0) return reject(Error('request-rejected'));
          resolve(value.data);
        } catch { reject(Error('invalid-response')); }
      });
    });
    req.setTimeout(3000, () => req.destroy(Error('request-timeout')));
    req.on('error', () => reject(Error('connection-failed')));
    if (data) req.write(data);
    req.end();
  });
}
function queueState(token) {
  return new Promise((resolve, reject) => {
    const url = new URL('/maasst/logs-ws', base);
    url.protocol = base.protocol === 'https:' ? 'wss:' : 'ws:';
    url.searchParams.set('token', token);
    const ws = new WebSocket(url);
    let finished = false;
    const timer = setTimeout(() => finish(null), 5000);
    function finish(value) {
      if (finished) return;
      finished = true;
      clearTimeout(timer);
      ws.terminate();
      if (value) resolve(value); else reject(Error('queue-state-unavailable'));
    }
    ws.on('message', bytes => {
      try {
        const message = JSON.parse(bytes.toString()).message || '';
        if (message.includes('没有正在执行或排队的任务')) finish('idle');
        else if (message.includes('当前任务正在执行') || message.includes('任务正在排队')) finish('busy');
      } catch { finish(null); }
    });
    ws.on('error', () => finish(null));
  });
}
function stable(value) {
  if (Array.isArray(value)) return value.map(stable);
  if (value && typeof value === 'object') {
    const fields = Object.keys(value).sort().filter(key => {
      if (key === 'lastExecutedAt') return false;
      // Existing jobs gain these nullable columns when the toolbox backend is deployed.
      return (key !== 'taskActivity' && key !== 'taskReclamation') || value[key] !== null;
    });
    return Object.fromEntries(fields.map(key => [key, stable(value[key])]));
  }
  return value;
}
const normalizedJobs = jobs => JSON.stringify(stable(jobs.slice().sort((a,b) => a.id-b.id)));
(async () => {
  if (!['snapshot','verify'].includes(mode) || !snapshotPath || !path.isAbsolute(snapshotPath)) throw Error('invalid-arguments');
  // Credentials stay local and are never included in stdout or error messages.
  const config = yaml.load(fs.readFileSync(path.join(root,'manifest/config/config.yaml'),'utf8'));
  const users = process.env.MAA_LOGIN_TOKENS ? Object.entries(JSON.parse(process.env.MAA_LOGIN_TOKENS)).map(([username,token]) => ({username,token})) : config.auth.users;
  if (!Array.isArray(users) || !users.length || users.length !== config.maa.asst.length) throw Error('incomplete-account-coverage');
  let sessions;
  if (mode === 'snapshot') {
    stage = 'check-login-schedules-and-idle-queues';
    sessions = [];
    for (const user of users) {
      const login = await request('POST','/maasst/login',{token:user.token});
      const jobs = await request('GET','/maasst/cronjob',undefined,login.token);
      const state = await queueState(login.token);
      sessions.push({username:user.username,token:login.token,jobs:jobs.cronJobList || [],queueState:state});
    }
    fs.mkdirSync(path.dirname(snapshotPath), {recursive:true});
    fs.writeFileSync(snapshotPath,JSON.stringify(sessions,null,2),{mode:0o600});
    const idle = sessions.every(s => s.queueState === 'idle');
    console.log(JSON.stringify({success:idle,accounts:sessions.length,queueStates:sessions.map(s=>s.queueState),scheduledJobCounts:sessions.map(s=>s.jobs.length)}));
    if (!idle) process.exitCode = 2;
  } else {
    stage = 'wait-for-backend';
    sessions = JSON.parse(fs.readFileSync(snapshotPath,'utf8'));
    if (sessions.length !== users.length || !sessions.length) throw Error('invalid-snapshot');
    let ready = false;
    for (let attempt=0; attempt<40; attempt++) {
      try { await request('GET','/maasst/cronjob',undefined,sessions[0].token); ready=true; break; } catch {}
      await new Promise(resolve=>setTimeout(resolve,1000));
    }
    if (!ready) throw Error('backend-not-ready');
    stage = 'verify-old-sessions-schedules-and-logins';
    for (const session of sessions) {
      const jobs = await request('GET','/maasst/cronjob',undefined,session.token);
      if (normalizedJobs(jobs.cronJobList||[]) !== normalizedJobs(session.jobs)) throw Error('schedule-data-changed');
      const user = users.find(u=>u.username===session.username);
      if (!user) throw Error('account-changed');
      await request('POST','/maasst/login',{token:user.token});
      await queueState(session.token);
    }
    console.log(JSON.stringify({success:true,oldSessionsAccepted:sessions.length,scheduleDataUnchanged:true,scheduledJobCounts:sessions.map(s=>s.jobs.length)}));
  }
})().catch(()=>{console.error('Backend maintenance check failed during '+stage+'; credentials withheld.');process.exitCode=1;});
