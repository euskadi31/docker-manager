import { Component, OnInit } from '@angular/core';

import { Network } from '../../entities/network';

import { NetworkService } from '../../services/network.service';

@Component({
    selector: 'app-network',
    templateUrl: './network.component.html',
    styleUrls: ['./network.component.css']
})
export class NetworkComponent implements OnInit {

    networks: Network[];

    constructor(private networkService: NetworkService) {
        this.networks = [];
    }

    ngOnInit() {
        this.networkService.getNetworks().then(networks => {
            console.log(networks);
            this.networks = networks;
        });
    }

}
