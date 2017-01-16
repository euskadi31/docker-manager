import { Version } from './version';

export class JoinTokens {
    Worker: string;
    Manager: string;
}

export class SwarmOrchestration {
    TaskHistoryRetentionLimit: number;
}

export class SwarmRaft {
    SnapshotInterval: number;
    LogEntriesForSlowFollowers: number;
    HeartbeatTick: number;
    ElectionTick: number;
}

export class SwarmDispatcher {
    HeartbeatPeriod: number;
}

export class SwarmCAConfig {
    NodeCertExpiry: number;
}

export class SwarmTaskDefaults {
}

export class SwarmSpec {
    Name: string;
    Orchestration: SwarmOrchestration;
    Raft: SwarmRaft;
    Dispatcher: SwarmDispatcher;
    CAConfig: SwarmCAConfig;
    TaskDefaults: SwarmTaskDefaults;
}

export class Swarm {
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    JoinTokens: JoinTokens;
    Spec: SwarmSpec;
}
