export async function atLeastOnePromise<T>(promises: Promise<T>[]): Promise<[T[], [number, Error][]]> {
    const results = await Promise.allSettled(promises);

    let fulfilled: T[] = [];
    let rejected: [number, Error][] = [];

    for (const [index, result] of results.entries()) {
        if (result.status === 'fulfilled') fulfilled.push(result.value);
        else rejected.push([index, result.reason]);
    };

    if (fulfilled.length === 0 && results.length > 0) {
        throw "All promises failed";
    }

    return [fulfilled, rejected];
}

export async function deepCopy<T>(obj: T): Promise<T> {
    return JSON.parse(JSON.stringify(obj));
}

export function deepEquality<T>(a: T, b: T): boolean {
    return JSON.stringify(a) === JSON.stringify(b);
}

export function isDescendentOf(descendent: HTMLElement, element: HTMLElement): boolean {
    for (let node: (HTMLElement | null) = descendent; node; node = node.parentElement) {
        if (node === element) return true;
    }
    return false;
}