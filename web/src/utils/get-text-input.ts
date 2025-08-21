

export function getInputWidth(value: string, font: string = "16px Arial") {
  // create virtual span
  const span = document.createElement("span");
  span.style.visibility = "hidden";
  span.style.position = "absolute";
  span.style.whiteSpace = "pre";
  span.style.font = font;
  span.textContent = value || " "; // use a space if empty
  document.body.appendChild(span);

  const width = span.offsetWidth + 8; // add some padding
  document.body.removeChild(span);

  return Math.max(width, 32); // minimum width
}