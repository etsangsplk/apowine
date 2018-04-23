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
    background-color: brown;
    color: white;
    font-size: 20px;
    font-family: arial;
  }

    #popupbox{
  position: relative;
  background: #FBFBF0;
  border: solid #000000 2px;
  font-family: arial;
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
    margin-top: 20%;
  }

  </style>
</head>

<body class="bodyStyle">
<h1 class="titleHeader">APOWINE</h1>
<div class="beer">
  <h2>BEER</h2>
  <h4>RANDOM BEER</h4>
  <button onclick="RandomDrink('beer');">Random</button>
</div>
<div class="wine">
  <h2>WINE</h2>
  <h4>RANDOM WINE</h4>
  <button onclick="RandomDrink('wine');">Random</button>
</div>
<div class="random">
  <h2>RANDOM DRINK</h2>
  <button onclick="RandomDrink('random');">Random Drink</button>
</div>
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
</body>
<script>
function CreateDrink(drinkType){

  var request = new XMLHttpRequest();
  if (drinkType == "beer"){
  var drinkName = document.getElementById("CbeerValue").value;
  request.open('POST', '/beer?type=create', true);
  request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  alert("BeerName: "+drinkName)
  request.send(JSON.stringify({beername:drinkName}));
}else{
  var drinkName = document.getElementById("CwineValue").value;
  request.open('POST', '/wine?type=create', true);
  request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  alert("WineName: "+drinkName)
  request.send(JSON.stringify({winename:drinkName}));
}
}

function ReadDrink(drinkType){
  var request = new XMLHttpRequest();
  if (drinkType=="beer"){
  var id = document.getElementById("RbeerID").value;
  request.open('GET', '/beer/'+id, true);
  request.onload=function(){
    var name = JSON.parse(request.response)
    alert("ID: "+ name.id+"\n"+ "BeerName: "+name.beername)
  }
  request.send();
}else{
  var id = document.getElementById("RwineID").value;
  request.open('GET', '/wine/'+id, true);
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
  request.open('PUT', '/beer', true);
  request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  request.send(JSON.stringify({id:drinkID,beername:drinkName}));
  alert("ID: "+ drinkID+"\n"+"BeerName: "+drinkName)
}else{
  var drinkID = document.getElementById("UwineID").value;
  var drinkName = document.getElementById("UwineValue").value;
  request.open('PUT', '/wine', true);
  request.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  request.send(JSON.stringify({id:drinkID,winename:drinkName}));
  alert("ID: "+ drinkID+"\n"+"WineName: "+drinkName)
}
}

function DeleteDrink(drinkType){
  var request = new XMLHttpRequest();
  if (drinkType=="beer"){
  var id = document.getElementById("DbeerID").value;
  request.open('DELETE', '/beer/'+id, true);
  request.send();
}else{
  var id = document.getElementById("DwineID").value;
  request.open('DELETE', '/wine/'+id, true);
  request.send();
}
}

function RandomDrink(drinkType){
  var request = new XMLHttpRequest();
  if (drinkType=="beer"){
  request.open('GET', '/beer?type=random', true);
  request.onload=function(){
    var name = JSON.parse(request.response)
    alert("ID: "+ name.id+"\n"+ "BeerName: "+name.beername)
  }
  request.send();
}else if (drinkType=="wine"){
  request.open('GET', '/wine?type=random', true);
  request.onload=function(){
    var name = JSON.parse(request.response)
    alert("ID: "+ name.id+"\n"+"WineName: "+name.winename)
  }
  request.send();
}else {
  request.open('GET', '/random?type=random', true);
  request.onload=function(){
    var name = JSON.parse(request.response)
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
</html>

`
