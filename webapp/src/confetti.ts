/**
 * License MIT, copyright Daniel Lundin 2019
 * https://github.com/daniel-lundin/dom-confetti
 *
 * License MIT, Copyright (c) 2021 Matteus Hemström
 *
 * Changes:
 *  - 2021-07-01: add modifyElement option
 *  - 2021-07-02: convert to typescript
 *  - 2022-06-08: add createElement option
 *  - 2022-06-08: set pointerEvents to none
 */

type CreateElementFn = () => HTMLElement;
type ModifyElementFn = (element: HTMLElement) => any;
type RandomFn = () => number;

interface Physics {
    x: number,
    y: number,
    z: number,
    wobble: number,
    wobbleSpeed: number,
    velocity: number,
    angle2D: number,
    angle3D: number,
    tiltAngle: number,
    tiltAngleSpeed: number,
}

interface Fetti {
  element: HTMLElement,
  physics: Physics,
}

interface ConfettiConfig {
    angle?: number;
    spread?: number;
    width?: string;
    height?: string;
    perspective?: string;
    duration?: number;
    dragFriction?: number;
    stagger?: number;
    startVelocity?: number;
    elementCount?: number;
    decay?: number;
    colors?: string[];
    random?: () => number;
    modifyElement?: (element: HTMLElement) => any;
    createElement?: () => HTMLElement;
}


const defaultColors = ["#a864fd", "#29cdff", "#78ff44", "#ff718d", "#fdff6a"];

function createElements(root: HTMLElement, elementCount: number, colors: string[], width: string, height: string, createElement: CreateElementFn | undefined, modifyElement: ModifyElementFn|undefined) {
  return Array.from({ length: elementCount }).map((_, index) => {
    const element = createElement ? createElement() : document.createElement("div");
    const color = colors[index % colors.length];
    element.style.backgroundColor = color;
    element.style.width = width;
    element.style.height = height;
    element.style.position = "absolute";
    element.style.willChange = "transform, opacity";
    element.style.opacity = "0"
    element.style.pointerEvents = "none";
    if (typeof modifyElement === 'function') {
      modifyElement(element);
    }
    root.appendChild(element);
    return element;
  });
}

function randomPhysics(angle: number, spread: number, startVelocity: number, random: RandomFn): Physics {
  const radAngle = angle * (Math.PI / 180);
  const radSpread = spread * (Math.PI / 180);
  return {
    x: 0,
    y: 0,
    z: 0,
    wobble: random() * 10,
    wobbleSpeed: 0.1 + random() * 0.1,
    velocity: startVelocity * 0.5 + random() * startVelocity,
    angle2D: -radAngle + (0.5 * radSpread - random() * radSpread),
    angle3D: -(Math.PI / 4) + random() * (Math.PI / 2),
    tiltAngle: random() * Math.PI,
    tiltAngleSpeed: 0.1 + random() * 0.3
  };
}

function updateFetti(fetti: Fetti, progress: number, dragFriction: number, decay: number) {
  /* eslint-disable no-param-reassign */
  fetti.physics.x += Math.cos(fetti.physics.angle2D) * fetti.physics.velocity;
  fetti.physics.y += Math.sin(fetti.physics.angle2D) * fetti.physics.velocity;
  fetti.physics.z += Math.sin(fetti.physics.angle3D) * fetti.physics.velocity;
  fetti.physics.wobble += fetti.physics.wobbleSpeed;
  // Backward compatibility
  if (decay) {
    fetti.physics.velocity *= decay;
  } else {
    fetti.physics.velocity -= fetti.physics.velocity * dragFriction;
  }
  fetti.physics.y += 3;
  fetti.physics.tiltAngle += fetti.physics.tiltAngleSpeed;

  const { x, y, z, tiltAngle, wobble } = fetti.physics;
  const wobbleX = x + 10 * Math.cos(wobble);
  const wobbleY = y + 10 * Math.sin(wobble);
  const transform = `translate3d(${wobbleX}px, ${wobbleY}px, ${z}px) rotate3d(1, 1, 1, ${tiltAngle}rad)`;

  fetti.element.style.transform = transform;
  fetti.element.style.opacity = (1 - progress).toString();

  /* eslint-enable */
}

function animate(root: HTMLElement, fettis: Fetti[], dragFriction: number, decay: number, duration: number, stagger: number) {
  let startTime = 0;

  return new Promise<void>(resolve => {
    function update(time: number) {
      if (startTime === 0) startTime = time;
      const elapsed = time - startTime;
      const progress = startTime === time ? 0 : (time - startTime) / duration;
      fettis.slice(0, Math.ceil(elapsed / stagger)).forEach(fetti => {
        updateFetti(fetti, progress, dragFriction, decay);
      });

      if (time - startTime < duration) {
        requestAnimationFrame(update);
      } else {
        fettis.forEach(fetti => {
          if (fetti.element.parentNode === root) {
            return root.removeChild(fetti.element);
          }
        });
        resolve();
      }
    }

    requestAnimationFrame(update);
  });
}

const defaults = {
  angle: 90,
  spread: 45,
  startVelocity: 45,
  elementCount: 50,
  width: "10px",
  height: "10px",
  perspective: "",
  colors: defaultColors,
  duration: 3000,
  stagger: 0,
  dragFriction: 0.1,
  random: Math.random,
  elementTextContent: null,
  modifyElement: null,
  createElement: null,
};

export function confetti(root: HTMLElement, config: ConfettiConfig = {}) {
  const {
    elementCount,
    colors,
    width,
    height,
    perspective,
    angle,
    spread,
    startVelocity,
    decay,
    dragFriction,
    duration,
    stagger,
    random,
    modifyElement,
    createElement,
  } = Object.assign({}, defaults, config);
  root.style.perspective = perspective;
  const elements = createElements(root, elementCount, colors, width, height, createElement, modifyElement);
  const fettis = elements.map(element => ({
    element,
    physics: randomPhysics(angle, spread, startVelocity, random)
  }));

  return animate(root, fettis, dragFriction, decay, duration, stagger);
}
