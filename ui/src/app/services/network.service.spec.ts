/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { NetworkService } from './network.service';

describe('NetworkService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [NetworkService]
    });
  });

  it('should ...', inject([NetworkService], (service: NetworkService) => {
    expect(service).toBeTruthy();
  }));
});
