export function capitalize(str: string) {
  return str[0]!.toUpperCase() + str.slice(1)
}

export function delay(ms: number) {
  return new Promise<void>((res) => {
    setTimeout(res, ms)
  })
}
