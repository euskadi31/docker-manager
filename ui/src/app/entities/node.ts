import { Version } from './version';

export class NodeSpec {
    Role: string;
    Availability: string;
}

export class NodePlatform {
    Architecture: string;
    OS: string;
}

export class NodeResources {
    NanoCPUs: number;
    MemoryBytes: number;
}

export class NodePlugin {
    Type: string;
    Name: string;
}

export class NodeEngine {
    EngineVersion: string;
    Plugins: NodePlugin[];
}

export class NodeStatus {
    State: string;
}

export class NodeManagerStatus {
    Leader: boolean;
    Reachability: string;
    Addr: string;
}

export class NodeDescription {
    Hostname: string;
    Platform: NodePlatform;
    Resources: NodeResources;
    Engine: NodeEngine;
    Status: NodeStatus;
    ManagerStatus: NodeManagerStatus;
}

export class Node {
    ID: string;
    Version: Version;
    CreatedAt: Date;
    UpdatedAt: Date;
    Spec: NodeSpec;
    Description: NodeDescription;
}
