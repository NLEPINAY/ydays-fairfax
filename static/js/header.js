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