import { TestBed } from '@angular/core/testing';

import { ElevatorStatsService } from './elevator-stats.service';

describe('ElevatorStatsService', () => {
  let service: ElevatorStatsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ElevatorStatsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
