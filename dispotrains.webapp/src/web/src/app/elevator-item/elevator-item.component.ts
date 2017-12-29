import { Component, Input } from '@angular/core';

import { Elevator } from '../station';

@Component({
  selector: 'elevator-item',
  templateUrl: './elevator-item.component.html',
  styleUrls: ['./elevator-item.component.css']
})
export class ElevatorItemComponent {
  @Input()
  elevator: Elevator;

  constructor() { }

  isBroken(): boolean {
    if (this.elevator === null) {
      return false;
    }
    return !this.elevator.available();
  }

  hasForecast(): boolean {
    if (this.elevator === null) {
      return false;
    }
    return this.elevator.status.forecast !== null;
  }
}
