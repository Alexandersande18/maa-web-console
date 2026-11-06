import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import test from 'node:test';
import { createDefaultTasks, getFieldOptions, isTaskFieldVisible } from '../src/components/task-config/taskDefinitions.js';
import { buildCronJobPayloadFromTasks, buildEnabledTasksPayload, processTaskConfig, syncTasksFromCronJob } from '../src/components/task-config/taskConfigUtils.js';

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
  assert.equal(tasks.taskRoguelike.config.starts_count, 3);
  assert.equal(tasks.taskRoguelike.config.blackflow_start_animals, '');
  assert.equal(buildEnabledTasksPayload(tasks, true).taskRoguelike.starts_count, 3);
  assert.equal(processTaskConfig('taskFight', { series: 11 }).series, 10);
  assert.equal(processTaskConfig('taskFight', { series: 0 }).series, 0);
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
