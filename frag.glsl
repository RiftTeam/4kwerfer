// 4k exegfx engine
// Rift, 2017

// This is an enhanced & simplified version of the renderer used for "TRSi Temple" and "Coming Up?". YOU should totally make an exegfx entry with it \o/

// Features:
//		- Ping pong float32 buffers in the basecode with automatic termination after max precalc time / max iters
//		- Pathtracer with diffuse/specular BRDF
//		- Classical RM used for scene traversal - make your distance functions the usual easy way!
//		- Unlimited (better than O(n)) emitters for your lighting pleasures
//		- Antialiasing, DOF, motion blur
//		- Easier cam orientation for exe gfx: cam pos, target pos, up
//		- Without the example scene it builds to 1899 in Fast, leaving over 2 full kb for your scene and texturers :)
//		- The alpha channel's completely unused, so you have a 32 bit channel to pass masks and stuff to the popro shader
//		- Debug builds automatically show progressive render; release builds show black until precalc time / max iters reached

// Todos (prod dependant so I haven't done them here):
//		- Colours-in-surface-definitions could be replaced with a GetTexture(p, n, surfaceID) function
//		- GetNormal func could take a surface ID too, for normal mapping based on surface
//		- Add another probability to the BRDF for refractive rays. Simple to do with sign-aware is() -- easiest to drop in our standard EST.
//		- Ditto for sss. Also simple to do with sign-aware is(); upon hitting the surface add its throughput, perturb and advance the ray, reflect it --> new p and rd.

#version 410

// ****************************** Constants ******************************

//////////////////////////////////////////////////////////////////////////////
// ENABLE FOR BUILD:
//#define PREVIOUS_FRAME c+=texture(fr,uv).xyz						// AUTOREP
// ENABLE FOR EDITOR:
#define PREVIOUS_FRAME c*=u.z+1										// AUTOREP
//////////////////////////////////////////////////////////////////////////////

#define SAMPLES_PER_FRAME 	6										// AUTOREP
#define MAX_BOUNCES 5												// AUTOREP
#define MAX_RM_STEPS 128											// AUTOREP
#define MAX_DIST 20													// AUTOREP
#define CAM_UP vec3(0,1,0)											// AUTOREP
#define CAM_POS vec3(0,-1,-6)										// AUTOREP
#define CAM_TARGET vec3(0,-.5,0)									// AUTOREP
#define FOV 5														// AUTOREP
#define FOCAL_LENGTH 6												// AUTOREP
//#define DOF_BLUR .05												// AUTOREP
#define DOF_BLUR 0													// AUTOREP
#define EPSILON_NORMAL .003											// AUTOREP
#define EPSILON_STEPOFF .015										// AUTOREP
#define EPSILON_RM .001												// AUTOREP
#define PI acos(-1)													// AUTOREP

// ****************************** Macros ******************************
#define x(q) r=q.x<r.x?q:r;											// Union
#define q(p,q) vec2(cos(q)*p.x+sin(q)*p.y,-sin(q)*p.x+cos(q)*p.y);	// 2D rotation of p by q

// ****************************** Globals ******************************
uniform vec3 u;														// {xres, yres, frame count}
uniform sampler2D fr;												// previous frame
out vec3 c;															// output col
vec2 uv;															// global uv
float fs;															// frame seed

// ****************************** Surface definitions ******************************
struct SF{
	float d,														// BRDF: Probability that a ray's diffuse (otherwise it's spec)
		g,															// glossiness
		e;															// emission: 1=emitter; 0=not emitter :)
	vec3 c;															// surface colour (RGB here - you could use HSV for better mixing, don't down-convert till the end)
};
SF sf[10]=SF[10](
	//brdf	gloss	emit	diffuse/emitted col
	SF(1,		0,		0,	vec3	(1,1,1)),								// white diffuse wall/ceiling
	SF(1,		0,		0,	vec3	(1,.6,.2)),								// orange diffuse wall
	SF(1,		0,		0,	vec3	(.2,1,.6)),								// green diffuse wall
	SF(0,		0,		1,	vec3	(.76,.84,.95)),							// light, emitter
	SF(1,		0,		0,	vec3	(1,1,1)),								// boxes
	SF(.1,	0,		0,	vec3	(.4,.2,1)),								// blue mirror ball
	SF(.1,	.35,	0,	vec3	(1,.2,.4)),								// red glossy ball
	SF(.5,	.5,		0,	vec3	(1,1,1)),								// glossy floor
	SF(0,		0,		1,	vec3	(.3,1,.2)),								// green alien technology
	SF(0,		0,		1,	vec3	(1,.2,.3))								// rudolph's glowing nose
);

// ****************************** Standard RM part ******************************
float vm(vec3 v){													// *** Get max vec comp, used by bx() ***
	return max(max(v.x,v.y),v.z);
}

float bx(vec3 p,vec3 b){											// *** Box prim ***
	vec3 d=abs(p)-b;
	return length(max(d,vec3(0)))+vm(min(d,vec3(0)));
}

vec2 h(vec3 p){														// *** Main world distance function; returns {dist, surface ID} ***
	vec2 r=vec2(MAX_DIST);
	x(vec2(bx(p-vec3(0,2,0),vec3(20,1,20)),0))						// ceiling
	x(vec2(bx(p-vec3(0,-3,0),vec3(20,1,20)),7))						// floor
	x(vec2(bx(p-vec3(0,0,3),vec3(20,20,1)),0))						// front wall
	x(vec2(bx(p-vec3(0,0,7),vec3(20,20,1)),0))						// back wall
	x(vec2(bx(p-vec3(-3,0,0),vec3(1,20,20)),1))						// left wall
	x(vec2(bx(p-vec3(3,0,0),vec3(1,20,20)),2))						// right wall
	
	vec3 b=p-vec3(1,-1.5,-.6);										// let's not mess with p for positioning stuff hmm? :) this is the elegant way to give this object group a local origin.
	b.xz=q(b.xz,-.5);
	x(vec2(bx(b,vec3(.5)),4))										// Right box
	x(vec2(length(b-vec3(0,.9,0))-.4,5))							// Ball
	
	b=p-vec3(-1,-2,-.3);											// fresh local origin
	b.xz=q(b.xz,.5);
	x(vec2(bx(b,vec3(.6)),4))										// left box
	x(vec2(length(b-vec3(0,1.1,0))-.5,6))							// ball
	
	x(vec2(bx(p-vec3(0,1.9,0),vec3(1)),3))							// main light
	x(vec2(length(p-vec3(.3,-1.9,-1))-.1,8))						// alien balls \o/
	x(vec2(length(p-vec3(-1.5,-1.9,-1-u.z*.0005))-.1,9))			// MOTION BLUR EXAMPLE: This alien ball is rolling towards us

	return r;
}

vec3 gn(vec3 p){													// *** Get normal at pos p for given surface ID ***
	vec2 e=vec2(EPSILON_NORMAL,0);
	return normalize(vec3(h(p+e.xyy).x-h(p-e.xyy).x,h(p+e.yxy).x-h(p-e.yxy).x,h(p+e.yyx).x-h(p-e.yyx).x));
}

vec2 is(vec3 p,vec3 r){												// *** Do classical RM to get intersection {dist, surface id} given origin and ray dir ***
	float t=0,d;
	for(int i=0;i<MAX_RM_STEPS;i++){
		d=h(p+t*r).x;
		if(abs(d)<EPSILON_RM||t>MAX_DIST)
			break;
		t+=d;
	}
	return vec2(t,h(p+t*r).y);
}

// *************************** Pathtracer part ***************************
float r1(vec3 p,float d){											// *** 1D rand given scale and seed ***
	return fract(sin(dot(gl_FragCoord.xyz+d,p))*43758.5453+d);		// yep, this stupid thing's where 100% of the engine's monte carlo choices come from, but we use some tricks when calling it to sample more cleverly.
}

vec2 r2(vec3 p){													// *** 3D->2D rand ***
	return vec2(r1(vec3(1),p.x+p.z),r1(vec3(1),p.y+p.z));			// just split the vec3 up and hit r1 twice
}

vec3 cd(vec3 n,float e){											// *** Cos-weighted sample centered around a given vector with given seed (based on http://www.rorydriscoll.com/2009/01/07/better-sampling/) ***
	float d=r1(vec3(12.9898,78.233,151.7182),e),					// you're spotting the pattern here right
		v=r1(vec3(63.7264,10.873,623.6736),e),
		r=sqrt(d),													// pick a random point on a disc - radius...
		a=PI*2*v;													// ...and theta
	vec3 p=abs(n.x)<.5?cross(n,vec3(1,0,0)):cross(n,vec3(0,1,0));	// compute basis from normal
	return r*cos(a)*p+r*sin(a)*cross(n,p)+sqrt(1-d)*n;				// and project it with a cosine distribution :)
}

vec3 pt(vec3 p,vec3 q){												// *** Calculate pixel colour for given ray origin and direction ***
    vec3 t=vec3(0),a=vec3(1);										// accumulated colour, throughput
	for(int i=0;i<MAX_BOUNCES;i++){									// we can't recurse so let's iterate the bounces :P
        vec2 r=is(p,q);												// check for an intersection along current raydir from current pos
		if(r.x>MAX_DIST)											// no hit?
			return t;												// return final col
		
		int m=int(r.y);												// surface/material ID
		vec3 h=p+q*r.x,												// hitpoint
			n=gn(h),												// normal
			s=sf[m].c;												// surface colour
	
		if(sf[m].e>0)												// hit an emitter?
			return i==0?s:t+a*s;									// if 1st iter, return the emitter col; otherwise, return final col (including this bounce)
			
		if(r1(vec3(1),uv.x+uv.y+fs+float(i))<sf[m].d)				// pick a ray type for next bounce based on surf type's BRDF probability (note: seed with everything possible)
			q=cd(n,fs*27.433727+i/MAX_BOUNCES);						// let's do a diffuse ray (seed from frame num and iters)
		else														// let's do a specular/reflective ray
			q=normalize(reflect(q,n)),								// compute a perfect one
			q+=cd(q,fs*9.312+i)*sf[m].g;							// perturb it by glossiness

		p=h+EPSILON_STEPOFF*q;										// set hitpoint to origin for next bounce
		a*=s;														// update throughput with this surface's colour
		t+=a*(1/MAX_BOUNCES);										// accumulate colour
    }
    return t;														// ah, we ran out of bounces without leaving scene or hitting an emitter. return what we got!
}

void main(){														// *** Entrypoint ***
	c=vec3(0);
	uv=gl_FragCoord.xy/u.xy;										// calc the actual pixel uv	
	for(int i=0;i<SAMPLES_PER_FRAME;i++){							// gonna do the whole thing multiple times \o/
		fs=u.z*(SAMPLES_PER_FRAME+1)+i;								// calculate seed for this frame		
		vec2 ns=r2(vec3(uv,fs)),									// get a noise val from curr uv and frame
			uv2=(-u.xy+2*(gl_FragCoord.xy+ns))/u.y;					// calc offset uv for AA
		vec3 ro=CAM_POS,											// ray origin
			ww=normalize(CAM_TARGET-CAM_POS),						// forward
			uu=normalize(cross(CAM_UP,ww)),							// right
			vv=normalize(cross(ww,uu)),								// it's polite to now recalc one's up
			er=normalize(vec3(uv2,FOV)),							// create a random ray w/DoF baby \o/ (Note, this is exactly iq's DOF code, http://www.iquilezles.org/www/articles/simplepathtracing/simplepathtracing.htm)
			rd=er.x*uu+er.y*vv+er.z*ww,								// calc rd -- note er uses FOV as z to force the xy comps to scale upon normalization
			go=DOF_BLUR*vec3(-1+2*ns,0),							// calc DOF lensing factors -- xy basis...
			gd=normalize(er*FOCAL_LENGTH-go);						// ...and z basis
		ro+=go.x*uu+go.y*vv;										// shift pos...
		rd+=gd.x*uu+gd.y*vv;										// ...and shift rd
		c+=pt(ro,normalize(rd));									// calc the pixel colour and accumulate it
	}
	c/=SAMPLES_PER_FRAME;											// scale by num samps we did
	PREVIOUS_FRAME;													// add to previous frame
  c=1-c;
}
// See? It fit in 200 lines <3
