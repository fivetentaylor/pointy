import { gql } from "@/__generated__";

export const GetMe = gql(`
  query GetMe {
    me {
      id
      email
      name
      displayName
      picture
      isAdmin

      subscriptionStatus
    }
  }
`);

export const GetMessagingLimit = gql(`
 query GetMessagingLimit {
    getMessagingLimits{
      used
      total
      startingAt
      endingAt
    }
  }
`);

export const GetUser = gql(`
  query GetUser($id: ID!) {
    user(id: $id) {
      id
      email
      name
      displayName
      picture
      isAdmin
    }
  }
`);

export const GetUsers = gql(`
  query GetUsers($ids: [ID!]!) {
    users(ids: $ids) {
      id
      email
      name
      displayName
      picture
      isAdmin
    }
  }
`);

export const getMyPreference = gql(`
  query GetMyPreference {
    myPreference {
      enableActivityNotifications
    }
  }
`);

export const updateMyPreference = gql(`
  mutation UpdateMyPreference($input: UpdateUserPreferenceInput!) {
    updateMyPreference(input: $input) {
      enableActivityNotifications
    }
  }
`);

export const updateMe = gql(`
  mutation UpdateMe($input: UpdateUserInput!) {
    updateMe(input: $input) {
      id
      email
      name
      displayName
      picture
      isAdmin
    }
  }
`);

export const GetUsersInMyDomain = gql(`
  query GetUsersInMyDomain($includeSelf: Boolean = false) {
    usersInMyDomain(includeSelf: $includeSelf) {
      id
      name
      displayName
      email
      picture
      isAdmin
    }
  }
`);
