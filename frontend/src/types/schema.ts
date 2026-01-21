import { MAX_PASSWORD_LEN } from '@/utils/constants'
import { boolean, object, type ObjectSchema, string } from 'yup'

export const loginSchema = object({
  username: string().required().min(3, 'Username must be at least 3 characters long'),
  password: string()
    .required()
    .max(MAX_PASSWORD_LEN, 'Password cannot be up to 50 characters in length'),
  stayLoggedIn: boolean(),
})

type SchemaValue<T> = T extends ObjectSchema<infer U> ? U : never
export type LoginSchema = SchemaValue<typeof loginSchema>
