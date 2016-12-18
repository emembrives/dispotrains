/* tslint:disable:no-unused-variable */
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

import { StationItemComponent } from './station-item.component';

describe('StationItemComponent', () => {
  let component: StationItemComponent;
  let fixture: ComponentFixture<StationItemComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StationItemComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StationItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
