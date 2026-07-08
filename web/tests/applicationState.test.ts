/*
 Pins the externally observable contract of the application state machine
 (components/stores/application-state.ts) after the xstate v4 -> v5 port, so a
 future refactor (or replacement of xstate entirely) can prove parity against
 this suite.

 Meta is asserted through metaForSnapshot(), the same function the app
 consumes it with (it flattens Object.values(snapshot.getMeta()) into one
 object).
*/

import { createActor, SimulatedClock, type AnyActorRef } from 'xstate';
import appStateModel, {
  AppStateEvent,
  metaForSnapshot,
} from '../components/stores/application-state';

const mergedMeta = metaForSnapshot;

// Expected merged meta per state, spelled out as field values (the meta
// constants for loading/goodbye are not exported from the module).
const LOADING_META = {
  chatAvailable: false,
  chatLoading: false,
  videoAvailable: false,
  appLoading: true,
};
const OFFLINE_META = {
  chatAvailable: false,
  chatLoading: false,
  videoAvailable: false,
  appLoading: false,
};
const ONLINE_META = {
  chatAvailable: true,
  chatLoading: false,
  videoAvailable: true,
  appLoading: false,
};
const GOODBYE_META = {
  chatAvailable: true,
  chatLoading: false,
  videoAvailable: false,
  appLoading: false,
};
const CHAT_USER_DISABLED_META = { ...ONLINE_META, chatAvailable: false };

const GOODBYE_TIMEOUT_MS = 300000;

const actors: AnyActorRef[] = [];
const startActor = (clock?: SimulatedClock) => {
  const actor = createActor(appStateModel, clock ? { clock } : undefined);
  actors.push(actor);
  actor.start();
  return actor;
};

afterEach(() => {
  // Stop every actor so a goodbye-state 'after' timer never outlives a test.
  actors.forEach(actor => actor.stop());
  actors.length = 0;
});

describe('application state machine', () => {
  test('starts in loading with the loading meta', () => {
    const actor = startActor();
    const snapshot = actor.getSnapshot();

    expect(snapshot.value).toBe('loading');
    expect(snapshot.matches('loading')).toBe(true);
    expect(snapshot.status).toBe('active');
    expect(mergedMeta(snapshot)).toEqual(LOADING_META);
  });

  test('NEEDS_REGISTER while loading stays in loading', () => {
    const actor = startActor();
    actor.send({ type: AppStateEvent.NeedsRegister });
    const snapshot = actor.getSnapshot();

    expect(snapshot.value).toBe('loading');
    expect(mergedMeta(snapshot)).toEqual(LOADING_META);
  });

  test('LOADED moves to ready.offline with the offline meta, and both matches() forms agree', () => {
    const actor = startActor();
    actor.send({ type: AppStateEvent.Loaded });
    const snapshot = actor.getSnapshot();

    expect(snapshot.value).toEqual({ ready: 'offline' });
    expect(snapshot.matches('ready')).toBe(true);
    expect(snapshot.matches({ ready: 'offline' })).toBe(true);
    expect(snapshot.matches('loading')).toBe(false);
    expect(mergedMeta(snapshot)).toEqual(OFFLINE_META);
  });

  test('full online walk: offline -> online -> goodbye -> back online', () => {
    const actor = startActor();
    actor.send({ type: AppStateEvent.Loaded });

    actor.send({ type: AppStateEvent.Online });
    let snapshot = actor.getSnapshot();
    expect(snapshot.matches({ ready: 'online' })).toBe(true);
    expect(mergedMeta(snapshot)).toEqual(ONLINE_META);

    actor.send({ type: AppStateEvent.Offline });
    snapshot = actor.getSnapshot();
    expect(snapshot.matches({ ready: 'goodbye' })).toBe(true);
    expect(snapshot.matches('ready')).toBe(true);
    expect(mergedMeta(snapshot)).toEqual(GOODBYE_META);

    actor.send({ type: AppStateEvent.Online });
    snapshot = actor.getSnapshot();
    expect(snapshot.matches({ ready: 'online' })).toBe(true);
    expect(mergedMeta(snapshot)).toEqual(ONLINE_META);
  });

  test('CHAT_USER_DISABLED while online keeps the online meta except chatAvailable', () => {
    const actor = startActor();
    actor.send({ type: AppStateEvent.Loaded });
    actor.send({ type: AppStateEvent.Online });
    actor.send({ type: AppStateEvent.ChatUserDisabled });
    const snapshot = actor.getSnapshot();

    expect(snapshot.matches({ ready: 'chatUserDisabled' })).toBe(true);
    expect(mergedMeta(snapshot)).toEqual(CHAT_USER_DISABLED_META);
  });

  test('FAIL while loading reaches serverFailure and the machine is done', () => {
    const actor = startActor();
    actor.send({ type: AppStateEvent.Fail });
    const snapshot = actor.getSnapshot();

    expect(snapshot.matches('serverFailure')).toBe(true);
    expect(snapshot.status).toBe('done');
  });

  test('ONLINE while loading is ignored without error', () => {
    const actor = startActor();
    actor.send({ type: AppStateEvent.Online });
    const snapshot = actor.getSnapshot();

    expect(snapshot.value).toBe('loading');
    expect(snapshot.status).toBe('active');
    expect(mergedMeta(snapshot)).toEqual(LOADING_META);
  });

  test('OFFLINE while ready.offline is ignored without error', () => {
    const actor = startActor();
    actor.send({ type: AppStateEvent.Loaded });
    actor.send({ type: AppStateEvent.Offline });
    const snapshot = actor.getSnapshot();

    expect(snapshot.matches({ ready: 'offline' })).toBe(true);
    expect(snapshot.status).toBe('active');
    expect(mergedMeta(snapshot)).toEqual(OFFLINE_META);
  });

  test('goodbye falls back to offline after exactly 300000ms', () => {
    const clock = new SimulatedClock();
    const actor = startActor(clock);
    actor.send({ type: AppStateEvent.Loaded });
    actor.send({ type: AppStateEvent.Online });
    actor.send({ type: AppStateEvent.Offline });
    expect(actor.getSnapshot().matches({ ready: 'goodbye' })).toBe(true);

    clock.increment(GOODBYE_TIMEOUT_MS - 1);
    expect(actor.getSnapshot().matches({ ready: 'goodbye' })).toBe(true);

    clock.increment(1);
    const snapshot = actor.getSnapshot();
    expect(snapshot.matches({ ready: 'offline' })).toBe(true);
    expect(mergedMeta(snapshot)).toEqual(OFFLINE_META);
  });

  test('leaving goodbye via ONLINE cancels the delayed fallback', () => {
    const clock = new SimulatedClock();
    const actor = startActor(clock);
    actor.send({ type: AppStateEvent.Loaded });
    actor.send({ type: AppStateEvent.Online });
    actor.send({ type: AppStateEvent.Offline });
    actor.send({ type: AppStateEvent.Online });

    clock.increment(GOODBYE_TIMEOUT_MS * 2);
    expect(actor.getSnapshot().matches({ ready: 'online' })).toBe(true);
  });
});
