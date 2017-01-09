import { Component, Input } from '@angular/core';
import { Location }         from '@angular/common';

@Component({
  selector: 'titlebar',
  templateUrl: './titlebar.component.html',
  styleUrls: ['./titlebar.component.css']
})
export class TitlebarComponent {
  @Input()
  title: string;
  @Input()
  root: boolean;

  constructor(private location: Location) { }

  goBack(): void {
    this.location.back();
  }

  hasBack(): boolean {
    return !this.root;
  }
}
