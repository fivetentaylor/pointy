import { gql } from "@/__generated__";

export const UploadImage = gql(`
  mutation UploadImage($file: Upload!, $docId: ID!) {
    uploadImage(file: $file, docId: $docId) {
      id
      docId
      url
      createdAt
      mimeType
      status
      error
    }
  }
`);

export const ListDocumentImages = gql(`
  query ListDocumentImages($docId: ID!) {
    listDocumentImages(docId: $docId) {
      id
      docId
      url
      createdAt
      mimeType
      status
      error
    }
  }
`);

export const GetImageSignedUrl = gql(`
  query GetImageSignedUrl($docId: ID!, $imageId: ID!) {
    getImageSignedUrl(docId: $docId, imageId: $imageId) {
      url
      expiresAt
    }
  }
`);

export const GetImage = gql(`
  query GetImage($docId: ID!, $imageId: ID!) {
    getImage(docId: $docId, imageId: $imageId) {
      id
      docId
      url
      createdAt
      mimeType
      status
      error
    }
  }
`);
