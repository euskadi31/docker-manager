import { Component, OnInit } from '@angular/core';

import { Image } from '../../entities/image';

import { ImageService } from '../../services/image.service';

@Component({
    selector: 'app-image',
    templateUrl: './image.component.html',
    styleUrls: ['./image.component.css']
})
export class ImageComponent implements OnInit {
    images: Image[];

    constructor(private imageService: ImageService) {
        this.images = [];
    }

    ngOnInit() {
        this.imageService.getImages().then(images => {
            console.log(images);
            this.images = images;
        });
    }

}
