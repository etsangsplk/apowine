package templates

// UITemplate to parse using template library
const UITemplate = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>APOWINE</title>
  <style>
  .bodyStyle{
    background-color: black;
    color: white;
    font-size: 20px;
    font-family: "Lucida Console", Monaco, monospace;
  }

    #popupbox{
  position: relative;
  font-family: "Lucida Console", Monaco, monospace;
  visibility: hidden;
  }

  .beer {
    float: left;
    margin-left: 250px;
    margin-top: 50px;
  }

  .wine {
    float: right;
    margin-right: 250px;
    margin-top: 50px;
  }

  .random{
    text-align: center;
    margin-top: 50px;
  }

  .titleHeader{
    text-align: center;
  }

  .more{
    text-align:center;
    margin-top: 25%;
  }

  .button {
      background-color: #008CBA;
      border: none;
      border-radius: 25px;
      color: white;
      padding: 100px 100px;
      text-align: center;
      text-decoration: none;
      display: inline-block;
      font-size: 20px;
      font-family: "Lucida Console", Monaco, monospace;
      margin: 40px 40px;
      cursor: pointer;
  }

  </style>
</head>

<body class="bodyStyle">
<h1 class="titleHeader">APOWINE</h1>
<div class="beer">
  <button class="button" onclick="RandomDrink('beer');">BEER</button>
</div>
<div class="wine">
  <button class="button" onclick="RandomDrink('wine');">WINE</button>
</div>
<div class="RbeerOP"><div>
<div class="more">
<div id="popupbox">
<form name="more" action="" method="post">

<div class="beer">
  <h2>BEER</h2>
  <h4>CREATE A BEER</h4>
  <input type="text" id="CbeerValue"  placeholder="beer name"/>
  <button onclick="CreateDrink('beer');">Create</button>
  <h4>READ BEER</h4>
  <input type="text" id="RbeerID" placeholder="beer ID"/>
  <button onclick="ReadDrink('beer');">Find</button>
  <h4>UPDATE A BEER</h4>
  <input type="text" id="UbeerID" placeholder="beer ID"/><br>
  <input type="text" id="UbeerValue" placeholder="beer name"/>
  <button onclick="UpdateDrink('beer');">Update</button>
  <h4>DELETE A BEER</h4>
  <input type="text" id="DbeerID" placeholder="beer ID"/>
  <button onclick="DeleteDrink('beer');">Delete</button>
</div>
<div class="wine">
  <h2>WINE</h2>
  <h4>CREATE A WINE</h4>
  <input type="text" id="CwineValue" placeholder="wine name"/>
  <button onclick="CreateDrink('wine');">Create</button>
  <h4>READ WINE</h4>
  <input type="text" id="RwineID" placeholder="wine ID"/>
  <button onclick="ReadDrink('wine');">Find</button>
  <h4>UPDATE A WINE</h4>
  <input type="text" id="UwineID" placeholder="wine ID"/><br>
  <input type="text" id="UwineValue" placeholder="wine name"/>
  <button onclick="UpdateDrink('wine');">Update</button>
  <h4>DELETE A WINE</h4>
  <input type="text" id="DwineID" placeholder="wine ID"/>
  <button onclick="DeleteDrink('wine');">Delete</button>
</div>

</form>
<br />
<center><a href="javascript:more('hide');">close</a></center>
</div>

<p><a href="javascript:more('show');">more>>></a></p>
</div>
<script type ="text/javascript">
function CreateDrink(drinkType){

  var request = new XMLHttpRequest();
  if (drinkType == "beer"){
  var drinkName = document.getElementById("CbeerValue").value;
  request.open('POST', '/drink?drinkType=beer&&operationType=create&&name='+drinkName, true);
  request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  request.onerror=function(e) {
    alert("Error creating object")
  }
  request.onload=function(e) {
    alert("Beer created")
  }
  request.send(JSON.stringify({beername:drinkName}));
}else{
  var drinkName = document.getElementById("CwineValue").value;
  request.open('POST', '/drink?drinkType=wine&&operationType=create&&name='+drinkName, true);
  request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  request.onerror=function(e) {
    alert("Error creating object")
  }
  request.onload=function(e) {
    alert("Wine created")
  }
  request.send(JSON.stringify({winename:drinkName}));
}
}

function ReadDrink(drinkType){
  var request = new XMLHttpRequest();
  if (drinkType=="beer"){
  var id = document.getElementById("RbeerID").value;
  request.open('GET', '/drink?drinkType=beer&&operationType=read&&id='+id, true);
  request.onload=function(){
    var name = JSON.parse(request.response)
    alert("ID: "+ name.id+"\n"+ "BeerName: "+name.beername)
  }
  request.send();
}else{
  var id = document.getElementById("RwineID").value;
  request.open('GET', '/drink?drinkType=wine&&operationType=read&&id='+id, true);
  request.onload=function(){
    var name = JSON.parse(request.response)
    alert("ID: "+ name.id+"\n"+"WineName: "+name.winename)
  }
  request.send();
}
}

function UpdateDrink(drinkType){
var request = new XMLHttpRequest();

if (drinkType == "beer"){
  var drinkID = document.getElementById("UbeerID").value;
  var drinkName = document.getElementById("UbeerValue").value;
  request.open('PUT', '/drink?drinkType=beer&&operationType=update&&id='+drinkID+'&&name='+drinkName, true);
  request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  request.send(JSON.stringify({id:drinkID,beername:drinkName}));
  alert("ID: "+ drinkID+"\n"+"BeerName: "+drinkName)
}else{
  var drinkID = document.getElementById("UwineID").value;
  var drinkName = document.getElementById("UwineValue").value;
  request.open('PUT', '/drink?drinkType=wine&&operationType=update&&id='+drinkID+'&&name='+drinkName, true);
  request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  request.send(JSON.stringify({id:drinkID,winename:drinkName}));
  alert("ID: "+ drinkID+"\n"+"WineName: "+drinkName)
}
}

function DeleteDrink(drinkType){
  var request = new XMLHttpRequest();
  if (drinkType=="beer"){
  var id = document.getElementById("DbeerID").value;
  request.open('DELETE', '/drink?drinkType=beer&&operationType=delete&&id='+id, true);
  request.send();
}else{
  var id = document.getElementById("DwineID").value;
  request.open('DELETE', '/drink?drinkType=wine&&operationType=delete&&id='+id, true);
  request.send();
}
}

function RandomDrink(drinkType){
  var request = new XMLHttpRequest();
  if (drinkType=="beer"){
  request.open('GET', '/drink?drinkType=beer&&operationType=random', true);
  request.onload=function(){
    var name = JSON.parse(request.response)
    alert("ID: "+ name.id+"\n"+ "BeerName: "+name.beername)
    document.getElementById("RbeerOP").innerHTML = "ID";
  }
  request.send();
}else if (drinkType=="wine"){
  request.open('GET', '/drink?drinkType=wine&&operationType=random', true);
  request.onload=function(){
    var name = JSON.parse(request.response)
    alert("ID: "+ name.id+"\n"+"WineName: "+name.winename)
  }
  request.send();
}else {
  request.open('GET', '/random', true);
  request.onload=function(){
    var name = JSON.parse(request.response)
    console.log(name)
    if (name.type=="beer"){
    alert("ID: "+ name.id+"\n"+ "BeerName: "+name.beername)
  }else{
    alert("ID: "+ name.id+"\n"+"WineName: "+name.winename)
  }
  }
  request.send();
}
}
function more(showhide){
if(showhide == "show"){
    document.getElementById('popupbox').style.visibility="visible";
}else if(showhide == "hide"){
    document.getElementById('popupbox').style.visibility="hidden";
}
}

</script>
</body>

</html>

`
