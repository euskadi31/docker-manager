import { Component, OnInit } from '@angular/core';

import { Node } from '../../entities/node';

import { NodeService } from '../../services/node.service';

@Component({
    selector: 'app-cluster',
    templateUrl: './cluster.component.html',
    styleUrls: ['./cluster.component.css']
})
export class ClusterComponent implements OnInit {

    nodes: Node[];

    constructor(private serviceService: NodeService) {
        this.nodes = [];
    }

    ngOnInit() {
        this.serviceService.getNodes().then(nodes => {
            console.log(nodes);
            this.nodes = nodes;
        });
    }

}
