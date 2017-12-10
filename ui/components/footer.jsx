/**
 * ui/components/footer.jsx
 * 
 * Defines the footer component to be placed at the bottom of the page.
 * 
 * @author zixian92
 */

'use strict'

import React, { Component } from 'react'

export default class Footer extends Component {
  render() {
    return (
      <div>
        Gringotts {VERSION}
      </div>
    )
  }
}
