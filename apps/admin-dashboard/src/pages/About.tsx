import { CheckCircle, Database, Shield, TrendingUp } from 'lucide-react';

export default function About() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-[#e8f5e9] to-white pb-20 md:pb-8">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          {/* Header */}
          <div className="text-center mb-12">
            <div className="inline-flex items-center justify-center bg-[#1a5f3f] text-white rounded-full p-4 mb-6">
              <TrendingUp className="w-12 h-12" />
            </div>
            <h1 className="text-4xl md:text-5xl font-bold text-[#1a5f3f] mb-4">
              About MarketLens
            </h1>
            <p className="text-xl text-gray-700">
              Agricultural Market Price Intelligence Platform
            </p>
          </div>
          
          {/* What We Do */}
          <div className="bg-white rounded-2xl shadow-lg p-8 mb-6">
            <div className="flex items-start gap-4 mb-4">
              <div className="bg-[#e8f5e9] rounded-lg p-3">
                <CheckCircle className="w-6 h-6 text-[#1a5f3f]" />
              </div>
              <div className="flex-1">
                <h2 className="text-2xl font-bold text-gray-900 mb-3">
                  What MarketLens Does
                </h2>
                <p className="text-gray-700 leading-relaxed mb-4">
                  MarketLens is a comprehensive price intelligence platform designed specifically 
                  for Nigeria's agricultural sector. We provide real-time market price information 
                  for crops across multiple markets, helping farmers, consumers, and retailers make 
                  informed decisions.
                </p>
                <p className="text-gray-700 leading-relaxed">
                  Our platform is not an e-commerce marketplace. We focus solely on providing 
                  accurate, timely price data to empower agricultural stakeholders with the 
                  information they need to maximize their profits and minimize losses.
                </p>
              </div>
            </div>
          </div>
          
          {/* How Prices Are Collected */}
          <div className="bg-white rounded-2xl shadow-lg p-8 mb-6">
            <div className="flex items-start gap-4 mb-4">
              <div className="bg-[#e8f5e9] rounded-lg p-3">
                <Database className="w-6 h-6 text-[#1a5f3f]" />
              </div>
              <div className="flex-1">
                <h2 className="text-2xl font-bold text-gray-900 mb-3">
                  How Prices Are Collected
                </h2>
                <p className="text-gray-700 leading-relaxed mb-4">
                  We gather price information through multiple trusted sources:
                </p>
                <ul className="space-y-3">
                  <li className="flex items-start gap-3">
                    <div className="bg-[#1a5f3f] rounded-full p-1 mt-1">
                      <CheckCircle className="w-4 h-4 text-white" />
                    </div>
                    <div>
                      <strong className="text-gray-900">Market Agents:</strong>
                      <span className="text-gray-700"> Trained field agents stationed at major markets who report daily prices</span>
                    </div>
                  </li>
                  <li className="flex items-start gap-3">
                    <div className="bg-[#1a5f3f] rounded-full p-1 mt-1">
                      <CheckCircle className="w-4 h-4 text-white" />
                    </div>
                    <div>
                      <strong className="text-gray-900">Mobile App Submissions:</strong>
                      <span className="text-gray-700"> Verified farmers and retailers can submit price observations</span>
                    </div>
                  </li>
                  <li className="flex items-start gap-3">
                    <div className="bg-[#1a5f3f] rounded-full p-1 mt-1">
                      <CheckCircle className="w-4 h-4 text-white" />
                    </div>
                    <div>
                      <strong className="text-gray-900">SMS Reports:</strong>
                      <span className="text-gray-700"> Network of contributors sending price updates via SMS</span>
                    </div>
                  </li>
                  <li className="flex items-start gap-3">
                    <div className="bg-[#1a5f3f] rounded-full p-1 mt-1">
                      <CheckCircle className="w-4 h-4 text-white" />
                    </div>
                    <div>
                      <strong className="text-gray-900">Government Data:</strong>
                      <span className="text-gray-700"> Official price bulletins from agricultural agencies</span>
                    </div>
                  </li>
                </ul>
                <p className="text-gray-700 leading-relaxed mt-4">
                  All submitted prices go through a verification process before being published 
                  on the platform. Our system aggregates multiple data points to provide the most 
                  accurate price ranges.
                </p>
              </div>
            </div>
          </div>
          
          {/* Confidence Levels */}
          <div className="bg-white rounded-2xl shadow-lg p-8 mb-6">
            <div className="flex items-start gap-4 mb-4">
              <div className="bg-[#e8f5e9] rounded-lg p-3">
                <Shield className="w-6 h-6 text-[#1a5f3f]" />
              </div>
              <div className="flex-1">
                <h2 className="text-2xl font-bold text-gray-900 mb-3">
                  Understanding Confidence Levels
                </h2>
                <p className="text-gray-700 leading-relaxed mb-4">
                  Each price listing includes a confidence indicator to help you assess data reliability:
                </p>
                <div className="space-y-4">
                  <div className="flex items-start gap-4 p-4 bg-green-50 rounded-lg border border-green-200">
                    <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium border border-green-200">
                      High
                    </span>
                    <div className="flex-1">
                      <p className="text-gray-700">
                        Multiple verified sources reporting consistent prices within the last 24 hours. 
                        This is the most reliable data.
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-start gap-4 p-4 bg-amber-50 rounded-lg border border-amber-200">
                    <span className="px-3 py-1 bg-amber-100 text-amber-700 rounded-full text-sm font-medium border border-amber-200">
                      Medium
                    </span>
                    <div className="flex-1">
                      <p className="text-gray-700">
                        Limited sources or slight variations in reported prices. Data is generally reliable 
                        but may need verification for large transactions.
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-start gap-4 p-4 bg-red-50 rounded-lg border border-red-200">
                    <span className="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm font-medium border border-red-200">
                      Low
                    </span>
                    <div className="flex-1">
                      <p className="text-gray-700">
                        Single source or outdated information. Use this data with caution and verify 
                        independently before making decisions.
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          {/* Disclaimer */}
          <div className="bg-amber-50 border-2 border-amber-300 rounded-2xl p-6">
            <h3 className="font-bold text-gray-900 mb-2 flex items-center gap-2">
              <Shield className="w-5 h-5 text-amber-600" />
              Important Disclaimer
            </h3>
            <p className="text-gray-700 leading-relaxed">
              Prices displayed on MarketLens are <strong>indicative only</strong> and updated daily. 
              Market conditions can change rapidly throughout the day. We strongly recommend verifying 
              current prices directly with markets before making significant purchase or sales decisions. 
              MarketLens is not responsible for any financial losses resulting from the use of this information.
            </p>
          </div>
          
          {/* Additional Info */}
          <div className="mt-6 text-center text-gray-600">
            <p className="mb-2">
              For questions or to become a price contributor, contact us at{' '}
              <a href="mailto:info@marketlens.ng" className="text-[#1a5f3f] hover:underline font-medium">
                info@marketlens.ng
              </a>
            </p>
            <p className="text-sm">
              Last updated: February 11, 2026
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
