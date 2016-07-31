#version 410
out vec3 outputColor;
uniform vec3 foo;
void main() {
    outputColor = vec3(.60,.9,foo.x);//texture(tex, fragTexCoord);
}
