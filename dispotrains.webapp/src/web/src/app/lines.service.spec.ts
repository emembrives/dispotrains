/* tslint:disable:no-unused-variable */

import { TestBed, async, inject } from '@angular/core/testing';
import { LinesService } from './lines.service';

describe('LinesService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [LinesService]
    });
  });

  it('should ...', inject([LinesService], (service: LinesService) => {
    expect(service).toBeTruthy();
  }));
});
