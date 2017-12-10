/**
 * ui/app.jsx
 * 
 * The entry point for UI code.
 * 
 * @author zixian92
 */

'use strict'

// Node packages
import React, { Component } from 'react'
import { render } from 'react-dom'

// Self-defined components
import Footer from './components/footer.jsx'

class App extends Component {
  render() {
    return (
      <div>
        <h1 className="cursor-pointer">Welcome to Gringotts!</h1>
        <Footer />
      </div>
    )
  }
}

render(<App />, document.getElementById('app'))
