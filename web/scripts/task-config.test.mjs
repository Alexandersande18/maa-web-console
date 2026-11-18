import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import test from 'node:test';
import { createDefaultTasks, getFieldOptions, getNumberFieldOptions, getReclamationStageGuide, isTaskFieldVisible } from '../src/components/task-config/taskDefinitions.js';
import { buildCronJobPayloadFromTasks, buildEnabledTasksPayload, normalizeNumberField, processTaskConfig, syncTasksFromCronJob } from '../src/components/task-config/taskConfigUtils.js';

test('relaunch anchor guidance matches the selected stage and is scoped to its theme', () => {
  for (const [mode, tokens, operator, steps] of [[16, 159, '无需指定干员', 1], [48, 417, '维什戴尔', 4], [32, 500, '圣聆初雪', 3]]) {
    const config = { theme: 'RelaunchAnchor', mode };
    const guide = getReclamationStageGuide(config);
    assert.match(guide[0].text, new RegExp(String(tokens)));
    assert.match(guide[1].text, new RegExp(operator));
    assert.equal(guide[2].steps.length, steps);
    assert.ok(guide[3].text.includes(getFieldOptions('taskReclamation', 'mode', config)[mode]));
  }
  assert.equal(getReclamationStageGuide({ theme: 'Tales', mode: 32 }), null);
  assert.equal(getReclamationStageGuide({ theme: 'RelaunchAnchor', mode: 0 }), null);
});

test('toolbox tasks preserve activity and manufacture options in a template', () => {
  const tasks = createDefaultTasks();
  tasks.taskActivity.enable = true;
  tasks.taskReclamation.enable = true;
  Object.assign(tasks.taskReclamation.config, { mode: 1, tools_to_craft: '荧光棒，发电机', clear_store: true });
  const saved = JSON.parse(JSON.stringify(buildCronJobPayloadFromTasks(tasks)));
  assert.equal(saved.taskActivity.activity_id, 'spa-regulations');
  assert.deepEqual(saved.taskReclamation.tools_to_craft, ['荧光棒', '发电机']);
  assert.equal(saved.taskReclamation.clear_store, undefined);
  const loaded = syncTasksFromCronJob(createDefaultTasks(), saved);
  assert.deepEqual(buildEnabledTasksPayload(loaded, true), {
    taskActivity: saved.taskActivity,
    taskReclamation: saved.taskReclamation
  });
  const anchor = { theme: 'RelaunchAnchor', mode: 48 };
  assert.equal(getFieldOptions('taskReclamation', 'mode', anchor)[48], 'RA-4');
  assert.equal(isTaskFieldVisible('taskReclamation', 'tools_to_craft', anchor), false);
  assert.throws(() => processTaskConfig('taskReclamation', { ...anchor, mode: 0 }), /主题和模式不匹配/);
});

test('BlackFlow configuration survives cron save, reload and immediate execution', () => {
  const tasks = createDefaultTasks();
  tasks.taskRoguelike.enable = true;
  Object.assign(tasks.taskRoguelike.config, {
    theme: 'BlackFlow', mode: 30003, difficulty: 6,
    blackflow_start_animals: 'swaddled_cat,swaddled_feathered_serpent',
    blackflow_start_rewards: 65, blackflow_start_utopia: true,
    collectible_mode_roles: '稳扎稳打',
    collectible_mode_core_char_list: '[{"name":"机械师","use_support":false}]'
  });
  const saved = buildCronJobPayloadFromTasks(tasks);
  assert.deepEqual(saved.taskRoguelike.blackflow_start_animals, ['swaddled_cat', 'swaddled_feathered_serpent']);
  assert.deepEqual(saved.taskRoguelike.collectible_mode_core_char_list, [{ name: '机械师', use_support: false }]);
  const loaded = syncTasksFromCronJob(createDefaultTasks(), JSON.parse(JSON.stringify(saved)));
  assert.deepEqual(buildEnabledTasksPayload(loaded, true).taskRoguelike, saved.taskRoguelike);
});

test('older cron jobs gain new fields and respect the exploration cap', () => {
  const tasks = syncTasksFromCronJob(createDefaultTasks(), {
    taskRoguelike: { enable: true, theme: 'JieGarden', starts_count: 8 }
  });
  assert.equal(tasks.taskRoguelike.config.starts_count, 8);
  assert.equal(tasks.taskRoguelike.config.blackflow_start_animals, '');
  assert.equal(buildEnabledTasksPayload(tasks, true).taskRoguelike.starts_count, 8);
  assert.equal(processTaskConfig('taskFight', { series: 11 }).series, 10);
  assert.equal(processTaskConfig('taskFight', { series: 0 }).series, 0);
});

test('selecting a saved template restores roguelike selection and JSON summary', () => {
  const tasks = createDefaultTasks();
  tasks.taskRoguelike.enable = true;
  Object.assign(tasks.taskRoguelike.config, { theme: 'BlackFlow', starts_count: 8, core_char: '机械师' });
  const template = JSON.parse(JSON.stringify({ isTemplate: 1, ...buildCronJobPayloadFromTasks(tasks) }));
  const loaded = createDefaultTasks();
  syncTasksFromCronJob(loaded, template);
  assert.equal(loaded.taskRoguelike.enable, true);
  assert.equal(loaded.taskRoguelike.config.theme, 'BlackFlow');
  const summary = JSON.parse(loaded.taskRoguelike.configSummary);
  assert.equal(summary.starts_count, 8);
  assert.equal(summary.core_char, '机械师');
  assert.deepEqual({ ...summary, enable: true }, buildEnabledTasksPayload(loaded, true).taskRoguelike);
});

test('switching templates clears null and missing tasks, selections and summaries', () => {
  const tasks = syncTasksFromCronJob(createDefaultTasks(), {
    taskRoguelike: { enable: true, theme: 'BlackFlow', starts_count: 8 },
    taskAward: { enable: true, award: false }
  });
  syncTasksFromCronJob(tasks, { taskRoguelike: null, taskFight: { enable: true, stage: '1-7' } });
  assert.deepEqual(Object.keys(buildEnabledTasksPayload(tasks, true)), ['taskFight']);
  const defaults = createDefaultTasks();
  assert.deepEqual(tasks.taskRoguelike, defaults.taskRoguelike);
  assert.deepEqual(tasks.taskAward, defaults.taskAward);
  syncTasksFromCronJob(tasks, { taskRoguelike: { enable: true, theme: 'Sami', starts_count: 3 } });
  assert.equal(tasks.taskRoguelike.enable, true);
  assert.equal(JSON.parse(tasks.taskRoguelike.configSummary).theme, 'Sami');
  assert.deepEqual(tasks.taskFight, defaults.taskFight);
});

test('exploration counts survive save, reload and immediate execution up to ten', () => {
  for (const [input, expected] of [[0, 0], [3, 3], [8, 8], [10, 10], [11, 10], [-1, 0], ['10', 10]]) {
    const tasks = createDefaultTasks();
    tasks.taskRoguelike.enable = true;
    tasks.taskRoguelike.config.starts_count = input;
    const saved = buildCronJobPayloadFromTasks(tasks);
    assert.equal(saved.taskRoguelike.starts_count, expected);
    const loaded = syncTasksFromCronJob(createDefaultTasks(), saved);
    assert.equal(loaded.taskRoguelike.config.starts_count, expected);
    assert.equal(buildEnabledTasksPayload(loaded, true).taskRoguelike.starts_count, expected);
  }
  assert.equal(createDefaultTasks().taskRoguelike.config.starts_count, 3);
});

test('numeric inputs preserve decimal thresholds, integer counts and sentinel values', () => {
  assert.equal(getNumberFieldOptions('taskInfrast', 'threshold').step, 0.1);
  assert.equal(normalizeNumberField('taskInfrast', 'threshold', '0.35'), 0.35);
  assert.equal(normalizeNumberField('taskInfrast', 'threshold', 2), 1);
  assert.equal(normalizeNumberField('taskFight', 'medicine', '-2'), 0);
  assert.equal(normalizeNumberField('taskRecruit', 'times', '4.8'), 4);
  assert.equal(normalizeNumberField('taskFight', 'times', -1), -1);
  assert.equal(normalizeNumberField('taskFight', 'times', '2147483647'), 2147483647);
  assert.equal(normalizeNumberField('taskRoguelike', 'difficulty', -1), -1);
  assert.equal(normalizeNumberField('taskRoguelike', 'blackflow_start_rewards', 5000), 4095);
  assert.equal(normalizeNumberField('taskRoguelike', 'starts_count', NaN), 3);
  assert.equal(normalizeNumberField('taskRoguelike', 'starts_count', Infinity), 3);
  const payload = processTaskConfig('taskFight', { times: '12', medicine: '3', series: '10' });
  assert.equal(payload.times, 12);
  assert.equal(payload.medicine, 3);
  assert.equal(payload.series, 10);
});

test('array parameters survive repeated serialization', () => {
  const config = processTaskConfig('taskRoguelike', {
    ...createDefaultTasks().taskRoguelike.config,
    start_foldartal_list: ['A', 'B'], expected_collapsal_paradigms: ['目空一些']
  });
  assert.deepEqual(processTaskConfig('taskRoguelike', config), config);
});

test('BlackFlow displays compatible modes and relevant parameters', () => {
  const config = { theme: 'BlackFlow', mode: 30003 };
  assert.deepEqual(Object.keys(getFieldOptions('taskRoguelike', 'mode', config)), ['0', '1', '30001', '30003']);
  assert.equal(isTaskFieldVisible('taskRoguelike', 'blackflow_start_animals', config), true);
  assert.equal(isTaskFieldVisible('taskRoguelike', 'blackflow_strategy', config), false);
  assert.equal(isTaskFieldVisible('taskRoguelike', 'blackflow_start_animals', { theme: 'Sami', mode: 5 }), false);
});

test('start-run rejects duplicate serpents and invalid Utopia settings', () => {
  const config = { ...createDefaultTasks().taskRoguelike.config, theme: 'BlackFlow', mode: 30003 };
  assert.throws(() => processTaskConfig('taskRoguelike', {
    ...config, blackflow_start_animals: 'swaddled_feathered_serpent,swaddled_feathered_serpent'
  }), /羽蛇/);
  assert.throws(() => processTaskConfig('taskRoguelike', {
    ...config, difficulty: 3, blackflow_start_utopia: true
  }), /难度6/);
});

test('standalone page and Vue app use identical task definitions', () => {
  for (const name of ['taskDefinitions.js', 'taskConfigUtils.js']) {
    assert.equal(readFileSync(new URL(`../src/components/task-config/${name}`, import.meta.url), 'utf8'),
      readFileSync(new URL(`../../resource/public/html/task-config/${name}`, import.meta.url), 'utf8'));
  }
});
