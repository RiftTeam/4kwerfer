#version 410
out vec3 outputColor;
in vec2 uv;
uniform vec3 u;
void main() {
//    if (uv.y > .5 + .5 *sin(u.z))//(gl_FragCoord.y > 200+ sin(gl_FragCoord.x/100)*50)
//    outputColor = vec3(.60,.9,foo.x + bar.y * baz);//texture(tex, fragTexCoord);
//    else outputColor = vec3(sin(uv.x*10 * 100*sin(u.z))*sin(uv.y*10 * u.z));
  outputColor = vec3(sin(u.z * 3));
}
