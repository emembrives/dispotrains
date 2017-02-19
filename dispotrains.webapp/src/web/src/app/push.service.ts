import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';

@Injectable()
export class PushService {
  private serverAPIBaseUrl = 'http://dispotrains.membrives.fr/app/SubscribePush/?endpoint=';

  constructor(private http: Http) {
  }

  registerPushAPI() {
    navigator.serviceWorker.register('push-sw.js', { scope: '/' }).then(this.registrationCallback).catch(function(error) {
      // registration failed
      console.log('Registration failed with ' + error);
    });;
  }

  private registrationCallback(registration: ServiceWorkerRegistration) {
    registration.pushManager.getSubscription().then((subscription: PushSubscription): Promise<PushSubscription> => {
      if (subscription) {
        return Promise.resolve(subscription);
      }
      return registration.pushManager.subscribe({ userVisibleOnly: false })
    }).then((subscription: PushSubscription) => {
      this.http.get(this.serverAPIBaseUrl + subscription.endpoint);
    }).catch(function(error) {
      // registration failed
      console.log('Subscription failed with ' + error);
    });
  }
}
