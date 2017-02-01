import { Injectable } from '@angular/core';

import { Image } from '../entities/image';

import { DockerService } from './docker.service';

@Injectable()
export class ImageService {
    constructor(private dockerService: DockerService) { }

    getImages(): Promise<Image[]> {
        return this.dockerService.get('/images/json')
            .then(response => response.json() as Image[])
            .then(images => {
                return images.map(image => {
                    if (image.RepoTags === null) {
                        [image.Repository] = image.RepoDigests[0].split('@')
                        image.Tag = '<none>';
                    } else {
                        [image.Repository, image.Tag] = image.RepoTags[0].split(':');
                    }

                    return image;
                })
            });
    }

    prune(): Promise<any> {
        return this.dockerService.post('/images/prune?filters={"dangling": true}', null)
            .then(response => response.json())
    }

    pull(name: string): Promise<any> {
        return this.dockerService.post(`/images/create?fromImage=${name}`, null);
    }
}
