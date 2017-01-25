import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs/Rx';

import { Service } from '../entities/service';
import { Log } from '../entities/log';

import { DockerService } from './docker.service';
import { WebSocketService } from './websocket.service';


@Injectable()
export class ServiceService {
    constructor(
        private dockerService: DockerService,
        private websocketService: WebSocketService
    ) {}

    getServices(): Promise<Service[]> {
        return this.dockerService.get("/services").then(response => response.json() as Service[]);
    }

    public getLogs(service: string): Subject<Log> {
        let protocol = (location.protocol === 'http:') ? 'ws:' : 'wss:';

        let url = `${protocol}//${location.host}/ws/service/${service}/log`;
        return <Subject<Log>>this.websocketService.connect(url)
            .map((response: MessageEvent): Log => JSON.parse(response.data) as Log);
    }

    getLogsOld(service: string): Observable<Log> {
        return new Observable(observer => {
            let protocol = (location.protocol === 'http:') ? 'ws:' : 'wss:';

            var ws = new WebSocket(`${protocol}//${location.host}/ws/service/${service}/log`);
            ws.addEventListener('error', () => observer.complete());
            ws.addEventListener('close', (event) => observer.complete())
            ws.addEventListener('message', (event) => observer.next(JSON.parse(event.data)), false);
        })
    }

    getLogsOld2(name): Promise<string[]> {
        return this.dockerService.get(`/services/${name}/logs?stdout=true`).then(response => {
            return response.text().split('\n')
        });
    }

    getLogsByTimestamp(name, since): Promise<string[]> {
        return this.dockerService.get(`/services/${name}/logs?stdout=true&since=${since}`).then(response => {
            return response.text().split('\n')
        });
    }
}
