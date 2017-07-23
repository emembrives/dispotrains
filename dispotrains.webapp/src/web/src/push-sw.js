self.addEventListener('push', function(event) {
  var notificationOptions = {
    body: "Hello World",
    icon: 'assets/logo-64.png',
    data:{
      url : 'http://example.com/updates'
    }
  };
  title = "Ceci est une notification !";
  return self.registration.showNotification(title, notificationOptions);
});
