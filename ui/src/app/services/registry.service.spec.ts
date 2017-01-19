/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { RegistryService } from './registry.service';

describe('RegistryService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [RegistryService]
    });
  });

  it('should ...', inject([RegistryService], (service: RegistryService) => {
    expect(service).toBeTruthy();
  }));
});
