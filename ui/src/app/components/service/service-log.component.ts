import { Component, OnInit, OnDestroy } from '@angular/core';
import {Router, ActivatedRoute, Params} from '@angular/router';

import { ServiceService } from '../../services/service.service';

@Component({
    selector: 'app-service-log',
    templateUrl: './service-log.component.html',
    styleUrls: ['./service-log.component.css']
})
export class ServiceLogComponent implements OnInit, OnDestroy {
    logs: string[];
    timer: NodeJS.Timer;
    since: number;

    constructor(
        private serviceService: ServiceService,
        private activatedRoute: ActivatedRoute
    ) {
        this.logs = [];
        this.since = 0;
    }

    fetchLog(name) {
        if (this.since === 0) {
            this.serviceService.getLogs(name).then(logs => {
                console.log(logs);
                this.logs = this.logs.concat(logs);
            });
        } else {
            this.serviceService.getLogsByTimestamp(name, this.since).then(logs => {
                console.log(logs);
                this.logs = this.logs.concat(logs);
            });
        }

        this.since = Math.floor(Date.now() / 1000);
    }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            this.fetchLog(params['name']);

            this.timer = setInterval(() => {
                this.fetchLog(params['name']);
            }, 1000);
        })
    }

    ngOnDestroy() {
        clearInterval(this.timer);
    }

}
