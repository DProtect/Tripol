'use strict';
importScripts('runtime.js');
self.onmessage=({data})=>{try{self.postMessage({ok:true,...TripolPlayground.run(data.source,data.input)})}catch(error){self.postMessage({ok:false,output:error.output||'',error:String(error.message||error)})}};
