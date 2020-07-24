//Load data on the header in every refresh of the component
$(document).ready(function() {

  var binumber = 0;
  for (i = 0; i < bichitos.length; i++){
      if (bichitos[i].status == "Online"){
        binumber++;
      }
  }

  var rednumber = 0;
  for (i = 0; i < redirectors.length; i++){
      rednumber++;
  }

  var inumber = 0;
  for (i = 0; i < implants.length; i++){
      inumber++;
  }

  var vnumber = 0;
  for (i = 0; i < vps.length; i++){
      vnumber++; 
  }

  var dnumber = 0;
  for (i = 0; i < domains.length; i++){
      dnumber++;
    
  }

  var snumber = 0;
  for (i = 0; i < stagings.length; i++){
      snumber++; 
  }
  
  $("#himplants").text(inumber);
  $("#hbichitos").text(binumber);
  $("#hdomains").text(dnumber);
  $("#hvps").text(vnumber);
  $("#hredirectors").text(rednumber);
  $("#hstagings").text(snumber);
  $("#username").text("Logged as "+ htmlencode.htmlEncode(username.username));

})
