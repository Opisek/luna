export const draggable = (node: HTMLElement) => {
  let moved = false;

  let similarElements: HTMLElement[] = [];
  let ownIndexInArray = 0;
  let lowestOrder = 0;

  let parent: HTMLElement;
  let phantom: HTMLElement;
  let wrapper: HTMLElement;

  let mouseOffsetY: number;
  let lowestY: number;
  let highestY: number;

  const mouseDown = () => {
    window.addEventListener("mouseup", mouseUp);
    window.addEventListener("mousemove", mouseMove);
    moved = false;
  }

  const mouseUp = (event: Event) => {
    window.removeEventListener("mouseup", mouseUp);
    window.removeEventListener("mousemove", mouseMove);

    if (moved) {
      // Prevent "click" events
      event.stopPropagation();
      event.stopImmediatePropagation();

      // Remove helper items
      phantom.remove();
      wrapper.remove();
      node.style.cursor = "";

      // Remove order attribute
      ([...parent.children] as HTMLElement[]).forEach(x => x.style.order = "");
      node.style.order = "";

      // Move the elements in the DOM
      if (ownIndexInArray == similarElements.length - 1) {
        if (lowestOrder + ownIndexInArray == parent.children.length - 1) {
          parent.appendChild(node);
        } else {
          parent.insertBefore(node, parent.children.item(lowestOrder + ownIndexInArray));
        }
      } else {
        parent.insertBefore(node, similarElements[ownIndexInArray + 1]);
      }
    }
  }

  const mouseMove = (event: Event) => {
    if (!node.parentElement) return;
    if (moved) {
      // Move the original element follow the mouse
      wrapper.style.top = `${Math.max(Math.min((event as MouseEvent).clientY - mouseOffsetY, highestY), lowestY)}px`;

      while (true) {
        const myBoundingBox = node.getBoundingClientRect();
        const myCenter = (myBoundingBox.top + myBoundingBox.bottom) / 2;

        let previousElement = ownIndexInArray == 0 ? null : similarElements[ownIndexInArray - 1];
        let nextElement = ownIndexInArray == similarElements.length - 1 ? null : similarElements[ownIndexInArray + 1];

        // Check if we moved the element up 
        if (previousElement) {
          const previousElementBoundingBox = previousElement.getBoundingClientRect();
          const previousElementCenter = (previousElementBoundingBox.top + previousElementBoundingBox.bottom) / 2;

          if (myCenter <= previousElementCenter + 1) {
            similarElements[ownIndexInArray] = previousElement;
            similarElements[ownIndexInArray - 1] = phantom;
            previousElement.style.order = (lowestOrder + ownIndexInArray).toString();
            phantom.style.order = (lowestOrder + ownIndexInArray - 1).toString();
            ownIndexInArray--;
          }
        }

        // Check if we moved the element down 
        if (nextElement) {
          const nextElementBoundingBox = nextElement.getBoundingClientRect();
          const nextElementCenter = (nextElementBoundingBox.top + nextElementBoundingBox.bottom) / 2;

          if (myCenter >= nextElementCenter - 1) {
            similarElements[ownIndexInArray] = nextElement;
            similarElements[ownIndexInArray + 1] = phantom;
            nextElement.style.order = (lowestOrder + ownIndexInArray).toString();
            phantom.style.order = (lowestOrder + ownIndexInArray + 1).toString();
            ownIndexInArray++;
          }
        }

        break;
      }
    } else {
      moved = true;

      parent = node.parentElement;

      // Find elements in own "ordering group" and add order property to every element in the container
      similarElements = [];
      ownIndexInArray = 0;
      let ownElementFound = false;
      let listComplete = false;
      for (const [i,x] of ([...parent.children] as HTMLElement[]).entries()) {
        x.style.order = i.toString();

        if (!listComplete) {
          if (x == node) {
            ownElementFound = true;
            ownIndexInArray = similarElements.length;
          }
          if (x.classList.contains("entry")) {
            if (similarElements.length == 0) lowestOrder = i;
            similarElements.push(x);
          }
          else if (ownElementFound) listComplete = true;
          else similarElements = [];
        }
      };

      if (!ownElementFound) return;

      if (similarElements.length == 1) {
        // Stop dragging if there is only one element to drag
        window.removeEventListener("mouseup", mouseUp);
        window.removeEventListener("mousemove", mouseMove);
        return;
      }

      const boundingRect = node.getBoundingClientRect();

      // Calculate where the element was pressed
      mouseOffsetY = (event as MouseEvent).clientY - boundingRect.top;
      lowestY = similarElements[0].getBoundingClientRect().top;
      highestY = similarElements[similarElements.length - 1].getBoundingClientRect().top;

      // Put a phantom element where the original was
      phantom = document.createElement("div");
      phantom.style.order = node.style.order;
      phantom.style.height = `${boundingRect.height}px`;
      phantom.style.width = `${boundingRect.width}px`;
      parent.replaceChild(phantom, node);
      similarElements[ownIndexInArray] = phantom;

      // Make the original element follow the mouse
      wrapper = document.createElement("div");
      wrapper.style.position = "fixed";
      wrapper.style.height = `${boundingRect.height}px`;
      wrapper.style.width = `${boundingRect.width}px`;
      wrapper.style.top = `${boundingRect.top}px`;
      wrapper.style.left = `${boundingRect.left}px`;
      document.body.appendChild(wrapper);
      wrapper.appendChild(node);

      // Change the cursor
      node.style.cursor = "ns-resize";

      // Prevent click events when we release the button later
      window.addEventListener("click", preventClick, true);
    }
  }

  const preventClick = (event: Event) => {
    event.stopPropagation();
    window.removeEventListener("click", preventClick, true);
  }

  // Detect when we start dragging
  node.addEventListener("mousedown", mouseDown);

  return {
    destroy() {
      node.removeEventListener("mousedown", mouseDown);
      window.removeEventListener("mouseup", mouseUp);
      window.removeEventListener("mousemove", mouseMove);
    }
  }
}