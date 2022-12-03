import { request } from '@/utils'

export function useUploadApi() {
  return { uploadFile: () => request.post('/upload') }
}
