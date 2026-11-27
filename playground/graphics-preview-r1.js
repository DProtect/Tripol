(function(root){
'use strict';
const rgb=v=>[(v>>>16&255)/255,(v>>>8&255)/255,(v&255)/255];
const valid=f=>f&&Number.isInteger(f.width)&&Number.isInteger(f.height)&&f.width>0&&f.height>0&&f.width<=512&&f.height<=512&&(f.kind==='rgb'&&f.rgba instanceof Uint8Array&&f.rgba.length===f.width*f.height*4||f.kind==='triangles'&&f.vertices instanceof Int32Array&&f.vertices.length%9===0&&f.vertices.length<=9000);
const pixelShader='@group(0) @binding(0) var img:texture_2d<f32>; @vertex fn vs(@builtin(vertex_index) i:u32)->@builtin(position) vec4f {var p=array<vec2f,3>(vec2f(-1.,-1.),vec2f(3.,-1.),vec2f(-1.,3.));return vec4f(p[i],0.,1.);} @fragment fn fs(@builtin(position) p:vec4f)->@location(0) vec4f{return textureLoad(img,vec2i(p.xy),0);}';
const triShader='struct V { @builtin(position) p:vec4f, @location(0) c:vec3f, }; @vertex fn vs(@location(0) p:vec2f,@location(1) c:vec3f)->V {var v:V;v.p=vec4f(p,0.,1.);v.c=c;return v;} @fragment fn fs(v:V)->@location(0) vec4f{return vec4f(v.c,1.);}';
async function gpuRenderer(canvas){
 if(!navigator.gpu)throw Error('WebGPU unavailable');
 const adapter=await navigator.gpu.requestAdapter();if(!adapter)throw Error('No WebGPU adapter');
 const device=await adapter.requestDevice(),context=canvas.getContext('webgpu');if(!context){device.destroy();throw Error('No WebGPU context')}
 const format=navigator.gpu.getPreferredCanvasFormat();context.configure({device,format,alphaMode:'opaque'});
 const pixel=device.createShaderModule({code:pixelShader}),tri=device.createShaderModule({code:triShader});
 device.pushErrorScope('validation');
 const pipeline=device.createRenderPipeline({layout:'auto',vertex:{module:pixel,entryPoint:'vs'},fragment:{module:pixel,entryPoint:'fs',targets:[{format}]},primitive:{topology:'triangle-list'}});
 const triangle=device.createRenderPipeline({layout:'auto',vertex:{module:tri,entryPoint:'vs',buffers:[{arrayStride:20,attributes:[{shaderLocation:0,offset:0,format:'float32x2'},{shaderLocation:1,offset:8,format:'float32x3'}]}]},fragment:{module:tri,entryPoint:'fs',targets:[{format}]},primitive:{topology:'triangle-list'}});
 const err=await device.popErrorScope();if(err){device.destroy();throw Error(err.message)}
 let lost=false;device.lost.then(()=>{lost=true});
 let texture=null,bind=null,textureWidth=0,textureHeight=0,vertexBuffer=null,vertexCapacity=0,vertexData=null;
 function release(resource){if(resource)device.queue.onSubmittedWorkDone().then(()=>resource.destroy(),()=>resource.destroy())}
 return {name:'WebGPU',draw(f){
  if(lost)throw Error('WebGPU device lost');
  if(canvas.width!==f.width)canvas.width=f.width;if(canvas.height!==f.height)canvas.height=f.height;
  const encoder=device.createCommandEncoder();let vertices=0;
  if(f.kind==='rgb'){
   if(textureWidth!==f.width||textureHeight!==f.height){
    release(texture);
    texture=device.createTexture({size:[f.width,f.height],format:'rgba8unorm',usage:GPUTextureUsage.TEXTURE_BINDING|GPUTextureUsage.COPY_DST});
    bind=device.createBindGroup({layout:pipeline.getBindGroupLayout(0),entries:[{binding:0,resource:texture.createView()}]});
    textureWidth=f.width;textureHeight=f.height;
   }
   device.queue.writeTexture({texture},f.rgba,{bytesPerRow:f.width*4},[f.width,f.height]);
  }else{
   vertices=f.vertices.length/3;
   if(vertices){
    const size=vertices*5;
    if(!vertexData||vertexData.length<size)vertexData=new Float32Array(size);
    for(let i=0,j=0;i<f.vertices.length;i+=3,j+=5){
     const color=f.vertices[i+2];vertexData[j]=f.vertices[i]/1000;vertexData[j+1]=f.vertices[i+1]/1000;
     vertexData[j+2]=(color>>>16&255)/255;vertexData[j+3]=(color>>>8&255)/255;vertexData[j+4]=(color&255)/255;
    }
    if(vertexCapacity<size*4){
     release(vertexBuffer);vertexCapacity=size*4;
     vertexBuffer=device.createBuffer({size:vertexCapacity,usage:GPUBufferUsage.VERTEX|GPUBufferUsage.COPY_DST});
    }
    device.queue.writeBuffer(vertexBuffer,0,vertexData,0,size);
   }
  }
  const [r,g,b]=rgb(f.background||0),pass=encoder.beginRenderPass({colorAttachments:[{view:context.getCurrentTexture().createView(),clearValue:{r,g,b,a:1},loadOp:'clear',storeOp:'store'}]});
  if(f.kind==='rgb'){pass.setPipeline(pipeline);pass.setBindGroup(0,bind);pass.draw(3)}
  else if(vertices){pass.setPipeline(triangle);pass.setVertexBuffer(0,vertexBuffer);pass.draw(vertices)}
  pass.end();device.queue.submit([encoder.finish()]);
 },dispose(){device.destroy()}};
}
function rasterTriangles(f,p=new Uint8ClampedArray(f.width*f.height*4)){
 const w=f.width,h=f.height,bg=rgb(f.background||0).map(x=>Math.round(x*255));for(let i=0;i<p.length;i+=4){p[i]=bg[0];p[i+1]=bg[1];p[i+2]=bg[2];p[i+3]=255}
 const edge=(a,b,x,y)=>(x-a.x)*(b.y-a.y)-(y-a.y)*(b.x-a.x);let work=0;
 for(let i=0;i<f.vertices.length;i+=9){
  const v=[];for(let j=0;j<9;j+=3)v.push({x:(f.vertices[i+j]/1000+1)*w/2,y:(1-f.vertices[i+j+1]/1000)*h/2,c:rgb(f.vertices[i+j+2])});
  const [a,b,c]=v,area=edge(a,b,c.x,c.y);if(Math.abs(area)<1e-8)continue;
  const x0=Math.max(0,Math.floor(Math.min(a.x,b.x,c.x))),x1=Math.min(w-1,Math.ceil(Math.max(a.x,b.x,c.x))),y0=Math.max(0,Math.floor(Math.min(a.y,b.y,c.y))),y1=Math.min(h-1,Math.ceil(Math.max(a.y,b.y,c.y)));
  work+=Math.max(0,x1-x0+1)*Math.max(0,y1-y0+1);if(work>4000000)throw Error('Software triangle preview limit');
  for(let y=y0;y<=y1;y++)for(let x=x0;x<=x1;x++){const u=edge(b,c,x+.5,y+.5)/area,vv=edge(c,a,x+.5,y+.5)/area,ww=1-u-vv;if(u>=-1e-8&&vv>=-1e-8&&ww>=-1e-8){const n=(y*w+x)*4;for(let k=0;k<3;k++)p[n+k]=255*(u*a.c[k]+vv*b.c[k]+ww*c.c[k])}}
 }return p;
}
function software(canvas){
 const ctx=canvas.getContext('2d');if(!ctx)throw Error('Canvas 2D unavailable');
 let image=null;
 return {name:'Canvas 2D',draw(f){
  if(!image||image.width!==f.width||image.height!==f.height){
   canvas.width=f.width;canvas.height=f.height;image=ctx.createImageData(f.width,f.height);
  }
  if(f.kind==='rgb')image.data.set(f.rgba);else rasterTriangles(f,image.data);
  ctx.putImageData(image,0,0);
 },dispose(){image=null}}
}
function create(host,language=()=>document.documentElement.lang){
 const say=(en,es)=>language()==='es'?es:en,heading=document.createElement('div'),label=document.createElement('span'),replay=document.createElement('button');heading.className='tripol-preview-head';replay.type='button';heading.append(label,replay);const stage=document.createElement('div');stage.className='tripol-preview-stage';let canvas=document.createElement('canvas');canvas.setAttribute('aria-label','TRIPOL graphics preview');stage.append(canvas);const status=document.createElement('p');status.className='tripol-preview-status';host.replaceChildren(heading,stage,status);host.hidden=true;
 let renderer=null,frames=[],timer=null,seq=0,disposed=false;
 const freshCanvas=()=>{const c=document.createElement('canvas');c.setAttribute('aria-label','TRIPOL graphics preview');canvas.replaceWith(c);canvas=c};
 function reset(){seq++;clearTimeout(timer);frames=[];host.hidden=true;if(renderer){renderer.dispose();renderer=null}freshCanvas()}
 async function play(){
  if(!frames.length||disposed)return;const id=++seq;clearTimeout(timer);host.hidden=false;replay.disabled=true;replay.textContent=say('Replay','Repetir');label.textContent=say('Graphics preview','Vista gráfica');status.textContent=say('Initializing renderer…','Iniciando renderizador…');
  if(!renderer){let pending;try{pending=await gpuRenderer(canvas)}catch{if(id!==seq||disposed)return;freshCanvas();try{pending=software(canvas)}catch(e){status.textContent=String(e.message);return}}if(id!==seq||disposed){pending.dispose();return}renderer=pending}
  let index=0;function step(){if(id!==seq||disposed)return;const f=frames[index];try{renderer.draw(f)}catch(e){if(renderer.name==='WebGPU'){renderer.dispose();freshCanvas();try{renderer=software(canvas);renderer.draw(f)}catch(err){status.textContent=String(err.message);replay.disabled=false;return}}else{status.textContent=String(e.message);replay.disabled=false;return}}
   label.textContent=(f.title||say('Graphics preview','Vista gráfica')).slice(0,128);status.textContent=renderer.name+' · '+f.width+' × '+f.height+' · '+(index+1)+' / '+frames.length;
   if(++index<frames.length)timer=setTimeout(step,Math.max(16,Math.min(1000,Number(f.delay)||16)));else replay.disabled=false
  }step();
 }
 replay.onclick=()=>play();
 return {reset,show(list){reset();frames=Array.isArray(list)?list.filter(valid).slice(0,120):[];if(frames.length)play()},dispose(){reset();disposed=true}};
}
root.TripolGraphics={create,rasterTriangles,valid};
})(globalThis);