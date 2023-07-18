uniform mat4 u_projectionMatrix;
attribute vec4 a_position;

void main() {
    vec4 pos = a_position;
    pos.xy *= 2.0;
    gl_position = u_projectionMatrix * a_position;
}

