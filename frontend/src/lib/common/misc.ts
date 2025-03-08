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