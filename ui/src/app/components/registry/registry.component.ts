import { Component, OnInit, ViewChild } from '@angular/core';

import { Registry } from '../../entities/registry';

import { RegistryService } from '../../services/registry.service';

import { DialogComponent } from '../dialog/dialog.component';

@Component({
    selector: 'app-registry',
    templateUrl: './registry.component.html',
    styleUrls: ['./registry.component.css']
})
export class RegistryComponent implements OnInit {

    registries: Registry[];

    registry: Registry;

    @ViewChild(DialogComponent)
    private dialog: DialogComponent;

    constructor(private registryService: RegistryService) {
        this.registries = [];
        this.registry = new Registry();
    }

    onFetchRegistries() {
        this.registryService.getRegistries().then(registries => {
            console.log(registries);
            this.registries = registries;
        });
    }

    ngOnInit() {
        this.onFetchRegistries();
    }

    onAdd() {
        this.registryService.addRegistry(this.registry).then(registry => {
            this.registry = new Registry();

            console.log('Registry:', registry);

            this.onFetchRegistries();

            this.dialog.close();
        });
    }

}
