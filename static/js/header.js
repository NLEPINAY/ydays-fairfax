
var params = new Object();
    params.action = "getStats";
    params.what = "Post";
    fetch("/fetching", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(params),
    })
      .then((x) => x.json()).then((x) => {console.log(x.DataChart[0].Count);
        if(x.DataChart[0].Count = 98){
            var elem =  document.getElementById("frind") ;
    
            if(getComputedStyle(elem).display != "none"){
                
                elem.style.display = "none";
            }else{
                elem.style.display = "block"; 
            }
        }
    
    });
 
      
  
function friend(){
    var elem =  document.getElementById("frind") ;
    
    if(getComputedStyle(elem).display != "none"){
        
        elem.style.display = "none";
    }else{
        elem.style.display = "block"; 
    }


}
function message(){
    var elem =  document.getElementById("mail") ;
    
    if(getComputedStyle(elem).display != "none"){
        
        elem.style.display = "none";
    }else{
        elem.style.display = "block"; 
    }


}