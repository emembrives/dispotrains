import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ElevatorStatsComponent } from './elevator-stats.component';

describe('ElevatorStatsComponent', () => {
  let component: ElevatorStatsComponent;
  let fixture: ComponentFixture<ElevatorStatsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [ElevatorStatsComponent]
    });
    fixture = TestBed.createComponent(ElevatorStatsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
