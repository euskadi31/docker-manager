import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptionsArgs, Response } from '@angular/http';
import 'rxjs/add/operator/toPromise';

const serialize = (obj: any) => {
    let str: any[] = [];
    for (let p in obj) {
        if (obj.hasOwnProperty(p)) {
            str.push(`${encodeURIComponent(p)}=${encodeURIComponent(obj[p])}`);
        }
    }

    return str.join('&');
};


@Injectable()
export class DockerService {
    private base = '/api';

    constructor(private http: Http) { }

    processBody(body: any, options: RequestOptionsArgs): any {
        if (options.headers.get('Content-Type') === 'application/x-www-form-urlencoded') {
            return serialize(body);
        }

        return body;
    }

    getOptions(options?: RequestOptionsArgs): RequestOptionsArgs {
        let headers = new Headers({
            'Content-Type': 'application/json',
            //'Authorization': this.getAuthorization(),
        });

        if (!options) {
            options = {};
        }

        if (options.headers) {
            headers.forEach((values, name) => {
                if (!options.headers.has(name)) {
                    options.headers.set(name, values);
                }
            });
        } else {
            options.headers = headers;
        }

        return options;
    }

    getUrl(path: string): string {
         if (path.charAt(0) === '/') {
             path = path.substr(1, path.length);
         }

         return `${this.base}/${path}`;
    }

    /**
     * Performs a request with `get` http method.
     */
    get(path: string, options?: RequestOptionsArgs): Promise<Response> {
        return this.http.get(this.getUrl(path), this.getOptions(options)).toPromise();
    }

    /**
     * Performs a request with `post` http method.
     */
    post(path: string, body: any, options?: RequestOptionsArgs): Promise<Response> {
        options = this.getOptions(options);

        return this.http.post(this.getUrl(path), this.processBody(body, options), options).toPromise();
    }

    /**
     * Performs a request with `put` http method.
     */
    put(path: string, body: any, options?: RequestOptionsArgs): Promise<Response> {
        options = this.getOptions(options);

        return this.http.put(this.getUrl(path), this.processBody(body, options), options).toPromise();
    }

    /**
     * Performs a request with `delete` http method.
     */
    delete(path: string, options?: RequestOptionsArgs): Promise<Response> {
        return this.http.delete(this.getUrl(path), this.getOptions(options)).toPromise();
    }
}
