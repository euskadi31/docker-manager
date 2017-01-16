import { Component, OnInit } from '@angular/core';

import { Service } from '../../entities/service';

import { ServiceService } from '../../services/service.service';

@Component({
    selector: 'app-service',
    templateUrl: './service.component.html',
    styleUrls: ['./service.component.css']
})
export class ServiceComponent implements OnInit {

    services: Service[];

    constructor(private serviceService: ServiceService) {
        this.services = [];
    }

    ngOnInit() {
        this.serviceService.getServices().then(services => {
            console.log(services);
            this.services = services;
        });
    }

}
