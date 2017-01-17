import { Component, OnInit } from '@angular/core';
import { PushService } from './push.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent implements OnInit {
  constructor(private pushService: PushService) { }

  ngOnInit(): void {
    this.pushService.registerPushAPI();
  }
}
