import { ComponentFixture, TestBed } from '@angular/core/testing';

import { StationItemComponent } from './station-item.component';

describe('StationItemComponent', () => {
  let component: StationItemComponent;
  let fixture: ComponentFixture<StationItemComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [StationItemComponent]
    });
    fixture = TestBed.createComponent(StationItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
