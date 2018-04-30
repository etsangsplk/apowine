function GetGoogleLoginPage(){
  var request = new XMLHttpRequest();
  request.open('GET', '/login?loginDomain=google', true);
  request.send();
}

function GetGithubLoginPage(){
  var request = new XMLHttpRequest();
  request.open('GET', '/login?loginDomain=github', true);
  request.send();
}
