import { Component, OnInit, OnDestroy } from '@angular/core';
import {Router, ActivatedRoute, Params} from '@angular/router';

import { Log } from '../../entities/log';

import { ServiceService } from '../../services/service.service';

@Component({
    selector: 'app-service-log',
    templateUrl: './service-log.component.html',
    styleUrls: ['./service-log.component.css']
})
export class ServiceLogComponent implements OnInit {
    logs: Log[];

    constructor(
        private serviceService: ServiceService,
        private activatedRoute: ActivatedRoute
    ) {
        this.logs = [];
    }

    ngOnInit() {
        this.activatedRoute.params.subscribe((params: Params) => {
            this.serviceService.getLogs(params['name']).subscribe(log => {
                console.log(log);
                this.logs.push(log);
            });
        })
    }
}
