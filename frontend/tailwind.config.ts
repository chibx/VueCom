import type { Config } from 'tailwindcss'
import colors from 'tailwindcss/colors'

export default <Partial<Config>>{
    theme: {
        extend: {
            colors: {
                primary: colors.green,
            },
        },
    },
    content: {
        files: [
            // all directories and extensions will correspond to your Nuxt config
            '{srcDir}/components/**/*.{vue,js,jsx,mjs,ts,tsx}',
            '{srcDir}/layouts/**/*.{vue,js,jsx,mjs,ts,tsx}',
            '{srcDir}/pages/**/*.{vue,js,jsx,mjs,ts,tsx}',
            '{srcDir}/{A,a}pp.{vue,js,jsx,mjs,ts,tsx}',
            '{srcDir}/{E,e}rror.{vue,js,jsx,mjs,ts,tsx}',
        ],
    },
}
