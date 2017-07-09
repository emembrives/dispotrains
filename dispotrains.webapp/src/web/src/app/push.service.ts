import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';
import { Observable }     from 'rxjs/Observable';
import { Subscription }     from 'rxjs/Subscription';

@Injectable()
export class PushService {
  private pushUrl = 'http://localhost:9000/app/PushSub';
  private vapidUrl = 'http://localhost:9000/app/VAPID';
  private registration: Subscription;

  constructor(private http: Http) {}

  private extractData(res: Response): string {
    let body = res.json();
    if (!body) {
      return "";
    }
    return body["PublicKey"];
  }

  private handleError(error: Response | any) {
    // In a real world app, we might use a remote logging infrastructure
    let errMsg: string;
    if (error instanceof Response) {
      const body = error.json() || '';
      const err = body.error || JSON.stringify(body);
      errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
    } else {
      errMsg = error.message ? error.message : error.toString();
    }
    console.error('Error while making network request: ' + errMsg);
    return Observable.throw(errMsg);
  }

  registerPushAPI() {
    this.registration = this.http.get(this.vapidUrl)
      .map(this.extractData)
      .catch(this.handleError)
      .map((key: string) => {
        navigator.serviceWorker.register('push-sw.js', { scope: '/' }).then((registration: ServiceWorkerRegistration) => {
          this.registrationCallback(key, registration);
        }).catch(function(error) {
          // registration failed
          console.log('Registration failed with ' + error);
        });
      }).subscribe();
  }

  private registrationCallback(key: string, registration: ServiceWorkerRegistration) {
    registration.pushManager.getSubscription().then((subscription: PushSubscription): Promise<PushSubscription> => {
      if (subscription) {
        return Promise.resolve(subscription);
      }
      let encodedKey = Uint8Array.from(Array.from(atob(key)).map((c : string) => {
        return c.codePointAt(0);
      }));
      return registration.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey: encodedKey,
      });
    }).then((subscription: PushSubscription) => {
      this.registration = this.http.post(this.pushUrl, subscription.toJSON())
        .subscribe();
    }).catch(function(error) {
      // registration failed
      console.log('Subscription failed with ' + error);
    });
  }
}
