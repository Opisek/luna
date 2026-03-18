import { effectiveBackgroundColor as getEffectiveBackground } from "../common/misc";

export const draggable = (node: HTMLElement, data: { ownClass: string, childClasses: string[], callback: (newIndex: number) => any }) => {
  let moved = false;

  let similarElements: HTMLElement[][] = [];
  let ownIndexInArray = 0;
  let originalIndexInArray = 0;

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

      // Move the elements in the DOM
      if (ownIndexInArray == similarElements.length - 1) {
        const myGroup = similarElements[ownIndexInArray];
        if (Number.parseInt(node.style.order) == parent.children.length - 1) {
          myGroup.forEach(x => {
            parent.appendChild(x);
          })
        } else {
          const nextElement = parent.children.item(Number.parseInt(node.style.order));
          similarElements[ownIndexInArray].forEach(x => {
            parent.insertBefore(x, nextElement);
          })
        }
      } else {
        const nextElement = similarElements[ownIndexInArray + 1][0];
        similarElements[ownIndexInArray].forEach(x => {
          parent.insertBefore(x, nextElement);
        })
      }

      // Remove order attribute
      ([...parent.children] as HTMLElement[]).forEach(x => x.style.order = "");

      if (ownIndexInArray != originalIndexInArray) data.callback(ownIndexInArray);
    }
  }

  const mouseMove = (event: Event) => {
    if (!node.parentElement) return;
    if (moved) {
      const myElementGroup = similarElements[ownIndexInArray];
      const myTop = myElementGroup[0].getBoundingClientRect().top;
      const myBottom = myElementGroup[myElementGroup.length - 1].getBoundingClientRect().bottom;
      const myHeight = myBottom - myTop;

      // Move the original element following the mouse
      wrapper.style.top = `${Math.max(Math.min((event as MouseEvent).clientY - mouseOffsetY, highestY - myHeight), lowestY)}px`;

      let anySwapped = false;
      let lowestSwapIndex = similarElements.length-1;
      let highestSwapIndex = 0;
      let lowestOriginalOrder = Number.MAX_SAFE_INTEGER;
      while (true) {
        const previousElementGroup = ownIndexInArray == 0 ? null : similarElements[ownIndexInArray - 1];
        const nextElementGroup = ownIndexInArray == similarElements.length - 1 ? null : similarElements[ownIndexInArray + 1];

        // Check if we moved the element up 
        if (previousElementGroup) {
          const previousElementGroupTop = previousElementGroup[0].getBoundingClientRect().top;
          const previousElementGroupBottom = previousElementGroup[previousElementGroup.length - 1].getBoundingClientRect().bottom;
          const previousElementGroupHeight = previousElementGroupBottom - previousElementGroupTop;

          const overlapRequired = previousElementGroupHeight / 3 * 2;

          if (myTop < previousElementGroupBottom - overlapRequired - 5) {
            anySwapped = true;
            lowestSwapIndex = Math.min(lowestSwapIndex, ownIndexInArray - 1);
            highestSwapIndex = Math.max(highestSwapIndex, ownIndexInArray);
            lowestOriginalOrder = Math.min(lowestOriginalOrder, Number.parseInt(previousElementGroup[0].style.order))

            similarElements[ownIndexInArray] = previousElementGroup;
            similarElements[ownIndexInArray - 1] = myElementGroup;
            ownIndexInArray--;
          }
        }

        // Check if we moved the element down 
        if (nextElementGroup) {
          const nextElementGroupTop = nextElementGroup[0].getBoundingClientRect().top;
          const nextElementGroupBottom = nextElementGroup[nextElementGroup.length - 1].getBoundingClientRect().bottom;
          const nextElementGroupHeight = nextElementGroupBottom - nextElementGroupTop;

          const overlapRequired = nextElementGroupHeight / 3 * 2;

          if (myBottom > nextElementGroupTop + overlapRequired + 5) {
            anySwapped = true;
            lowestSwapIndex = Math.min(lowestSwapIndex, ownIndexInArray);
            highestSwapIndex = Math.max(highestSwapIndex, ownIndexInArray + 1);
            lowestOriginalOrder = Math.min(lowestOriginalOrder, Number.parseInt(myElementGroup[0].style.order))

            similarElements[ownIndexInArray] = nextElementGroup;
            similarElements[ownIndexInArray + 1] = myElementGroup;
            ownIndexInArray++;
          }
        }

        break;
      }

      if (anySwapped) {
        let runningOrder = lowestOriginalOrder;
        for (let i = lowestSwapIndex; i <= highestSwapIndex; i++) {
          for (const element of similarElements[i]) {
            element.style.order = `${runningOrder++}`;
          }
        }
        phantom.style.order = node.style.order;
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
            originalIndexInArray = ownIndexInArray;
          }
          if (x.classList.contains(data.ownClass)) {
            similarElements.push([x]);
          } else if (similarElements.length != 0 && data.childClasses.some(c => x.classList.contains(c))) {
            similarElements[similarElements.length - 1].push(x);
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
      const myElementGroup = similarElements[ownIndexInArray];
      const lastChildBoundingRect = myElementGroup[myElementGroup.length - 1].getBoundingClientRect();

      // Calculate where the element was pressed
      mouseOffsetY = (event as MouseEvent).clientY - boundingRect.top;
      lowestY = similarElements[0][0].getBoundingClientRect().top;
      const lastElementGroup = similarElements[similarElements.length - 1]
      highestY = lastElementGroup[lastElementGroup.length - 1].getBoundingClientRect().bottom;

      // Make the original element following the mouse
      wrapper = document.createElement("div");
      wrapper.style.position = "fixed";
      wrapper.style.height = `${lastChildBoundingRect.bottom - boundingRect.top}px`;
      wrapper.style.width = `${boundingRect.width}px`;
      wrapper.style.top = `${boundingRect.top}px`;
      wrapper.style.left = `${boundingRect.left}px`;
      wrapper.style.display = "flex";
      wrapper.style.flexDirection = "column";
      wrapper.style.justifyContent = "space-between";
      wrapper.style.background = getEffectiveBackground(node);
      document.body.appendChild(wrapper);

      // Put a phantom element where the original was
      phantom = document.createElement("div");
      phantom.style.order = node.style.order;
      phantom.style.height = `${lastChildBoundingRect.bottom - boundingRect.top}px`;
      phantom.style.width = `${boundingRect.width}px`;
      parent.replaceChild(phantom, node);

      myElementGroup.forEach(x => {
        wrapper.appendChild(x);
      })

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