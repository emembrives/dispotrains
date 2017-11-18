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

  isBroken() : boolean {
    return !this.elevator.available();
  }

  hasForecast() : boolean {
    return this.elevator.status.forecast != undefined
  }
}
